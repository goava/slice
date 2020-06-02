package slice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goava/di"
)

const (
	// default start/shutdown timeout
	defaultTimeout = 5 * time.Second
)

// Run creates and runs application with default shutdown flow (SIGTERM, SIGINT).
func Run(options ...Option) {
	app := New(options...)
	if err := app.Start(); err != nil {
		app.logger.Fatal(err)
	}
}

// New creates slice application with provided options.
func New(options ...Option) *Application {
	s := Application{}
	for _, opt := range options {
		opt.apply(&s)
	}
	if s.configurator == nil {
		s.configurator = defaultBundleConfigurator()
	}
	if s.timeouts.start == 0 {
		s.timeouts.start = defaultTimeout
	}
	if s.timeouts.shutdown == 0 {
		s.timeouts.shutdown = defaultTimeout
	}
	return &s
}

// Application is a controlling part of application.
type Application struct {
	stop         func()
	di           []di.Option
	bundles      []Bundle
	logger       Logger
	dispatcher   Dispatcher
	configurator bundleConfigurator
	timeouts     struct {
		start    time.Duration
		shutdown time.Duration
	}
}

// Starts start slice.
func (app *Application) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), app.timeouts.start)
	// add application context
	app.di = append(app.di, di.Provide(func() context.Context { return ctx }))
	// save context cancel
	app.stop = cancel
	// sort bundles
	sorted, ok := sortBundles(app.bundles)
	if !ok {
		return fmt.Errorf("bundle cyclic detected") // todo: improve error message
	}
	// load bundle information
	bundles := inspectBundles(sorted...)
	// configure bundles
	if err := configureBundles(app.configurator, bundles...); err != nil {
		return err
	}
	// create dependency injection container
	container, err := createContainer(app.di...)
	if err != nil {
		return err
	}
	// build bundle dependencies
	if err := buildBundles(container, bundles...); err != nil {
		return err
	}
	// resolve logger
	_ = container.Resolve(&app.logger)
	// if logger not provided use std logger
	if app.logger == nil {
		app.logger = &stdLogger{}
		// provide std logger to container
		_ = container.Provide(func() Logger { return app.logger })
	}
	// run goroutine with os exit listener
	go app.waitExit()
	// boot bundles
	shutdowns, err := boot(ctx, container, bundles...)
	// if boot failed shutdown booted bundles
	if err != nil {
		if rserr := reverseShutdown(app.timeouts.shutdown, container, shutdowns); rserr != nil {
			return fmt.Errorf("%w (%s)", err, rserr)
		}
		return err
	}
	app.logger.Info("Start")
	// run application, ignore context cancel error
	// default context lifecycle used for application shutdown
	if err := run(ctx, container); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	app.logger.Info("Shutdown")
	// shutdown bundles in reverse order
	if err = reverseShutdown(app.timeouts.shutdown, container, shutdowns); err != nil {
		return fmt.Errorf("%w", err)
	}
	return err
}

func (app *Application) waitExit() {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	app.stop()
}

// Stop stops application.
func (app *Application) Stop() {
	app.stop()
}
