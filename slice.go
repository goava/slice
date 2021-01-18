package slice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/goava/di"
)

const (
	// default start/shutdown timeout
	defaultTimeout = 5 * time.Second
	defaultEnv     = "ENV"
	defaultDebug   = "DEBUG"
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
	s := Application{
		info: &Info{},
	}
	for _, opt := range options {
		opt.apply(&s)
	}
	env, _ := lookupEnv(defaultEnv)
	s.info.Env = parseEnv(env)
	if debug, ok := lookupEnv(defaultDebug); ok {
		s.info.Debug = strings.ToLower(debug) == "true"
	}
	if s.configurator == nil {
		s.configurator = defaultBundleConfigurator()
	}
	if s.timeouts.boot == 0 {
		s.timeouts.boot = defaultTimeout
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
		boot     time.Duration
		shutdown time.Duration
	}
	info *Info
}

// Starts start slice.
func (app *Application) Start() error {
	app.logger = &stdLogger{}
	if app.info == nil || app.info.Name == "" {
		return fmt.Errorf("application name must be specified, see slice.SetName() option")
	}
	base, stop := context.WithCancel(context.Background())
	ctx := NewContext(base)
	// store context cancel
	app.stop = stop
	// add application context
	app.di = append(app.di, di.Provide(func() *Context { return ctx }, di.As(new(context.Context))))
	app.di = append(app.di, di.Provide(func() Info { return *app.info }))
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
	// resolve logger from container
	// if logger not found it will remain std
	_ = container.Resolve(&app.logger)
	// start goroutine with os signal catch
	go app.catchSignals()
	startCtx, _ := context.WithTimeout(ctx, app.timeouts.boot)
	// boot bundles
	shutdowns, err := boot(startCtx, container, bundles...)
	// if boot failed shutdown booted bundles
	if err != nil {
		// create context for shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), app.timeouts.shutdown)
		defer cancel()
		if rserr := reverseShutdown(shutdownCtx, container, shutdowns); rserr != nil {
			return fmt.Errorf("%w (%s)", err, rserr)
		}
		return err
	}
	app.logger.Printf("Start")
	// run application, ignore context cancel error
	// default context lifecycle used for application shutdown
	if err := run(ctx, container); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	app.logger.Printf("Shutdown")
	// create context for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), app.timeouts.shutdown)
	defer cancel()
	// shutdown bundles in reverse order
	if err = reverseShutdown(shutdownCtx, container, shutdowns); err != nil {
		return fmt.Errorf("%w", err)
	}
	return err
}

// Stop stops application.
func (app *Application) Stop() {
	app.stop()
}

// catchSignals waits SIGTERM or SIGINT signals
func (app *Application) catchSignals() {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	app.logger.Printf(strings.Title(sign.String()))
	app.stop()
}
