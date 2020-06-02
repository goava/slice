package slice

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goava/di"
)

// Run creates and runs application with default shutdown flow (SIGTERM, SIGINT).
func Run(options ...Option) {
	lf := New(options...)
	if err := lf.Start(); err != nil {
		exitError(err)
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
		s.timeouts.start = 5 * time.Second
	}
	if s.timeouts.shutdown == 0 {
		s.timeouts.shutdown = 5 * time.Second
	}
	return &s
}

// Application is a controlling part of application.
type Application struct {
	stop         func()
	di           []di.Option
	bundles      []Bundle
	container    *di.Container
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
	app.di = append(app.di, di.Provide(func() context.Context { return ctx }))
	app.stop = cancel
	sorted, ok := sortBundles(app.bundles)
	if !ok {
		return fmt.Errorf("bundle cyclic detected") // todo: improve error message
	}
	bundles := inspectBundles(sorted...)
	if err := configureBundles(app.configurator, bundles...); err != nil {
		return err
	}
	container, err := createContainer(app.di...)
	if err != nil {
		return err
	}
	app.container = container
	if err := buildBundles(container, bundles...); err != nil {
		return err
	}
	shutdowns, err := boot(ctx, container, bundles...)
	if err != nil {
		// todo: timeout
		if rserr := reverseShutdown(app.timeouts.shutdown, container, shutdowns); rserr != nil {
			return fmt.Errorf("%s (%s)", err, rserr)
		}
		return err
	}
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-stop
		cancel()
	}()
	if err := run(ctx, container); err != nil {
		return err
	}
	if rserr := reverseShutdown(app.timeouts.shutdown, container, shutdowns); rserr != nil {
		return fmt.Errorf("%s (%s)", err, rserr)
	}
	return err
}

// Stop stops application.
func (app *Application) Stop() {
	app.stop()
}
