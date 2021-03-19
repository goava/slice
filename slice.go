package slice

import (
	"context"
	"errors"
	"flag"
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
		app.Logger.Fatal(err)
	}
}

// New creates slice application with provided options.
func New(options ...Option) *Application {
	s := Application{}
	for _, opt := range options {
		opt.apply(&s)
	}
	return &s
}

// Application is a controlling part of application.
type Application struct {
	Name            string
	Env             Env
	Debug           bool
	Prefix          string
	Parameters      []Parameter
	Components      []di.Option
	Dispatcher      Dispatcher
	Bundles         []Bundle
	StartTimeout    time.Duration
	StopTimeout     time.Duration
	Logger          Logger
	ParameterParser ParameterParser

	stop func()
}

// Starts start slice.
func (app *Application) Start() error {
	// set defaults
	if app.Logger == nil {
		app.Logger = &stdLogger{}
	}
	if app.ParameterParser == nil {
		app.ParameterParser = &stdParameterParser{}
	}
	if app.StartTimeout == 0 {
		app.StartTimeout = defaultTimeout
	}
	if app.StopTimeout == 0 {
		app.StopTimeout = defaultTimeout
	}
	// check application name
	if len(app.Name) == 0 {
		return fmt.Errorf("application name must be specified, see slice.SetName() option")
	}
	// lookup environment
	env, _ := lookupEnv(defaultEnv)
	app.Env = parseEnv(env)
	if debug, ok := lookupEnv(defaultDebug); ok {
		app.Debug = strings.ToLower(debug) == "true"
	}
	// initialize context
	base, stop := context.WithCancel(context.Background())
	ctx := NewContext(base)
	// store context cancel
	app.stop = stop
	info := Info{
		Name:  app.Name,
		Env:   app.Env,
		Debug: app.Debug,
	}
	// add application context and info
	components := append(app.Components,
		di.Provide(func() *Context { return ctx }, di.As(new(context.Context))),
		di.Provide(func() Info { return info }),
	)
	// sort bundles
	sorted, ok := sortBundles(app.Bundles)
	if !ok {
		return fmt.Errorf("bundle cyclic detected") // todo: improve error message
	}
	// create dependency injection container
	container, err := createContainer(components...)
	if err != nil {
		return err
	}
	parameters := app.Parameters
	// add bundle parameters
	for _, bundle := range sorted {
		parameters = append(parameters, bundle.Parameters...)
	}
	var help bool
	flag.BoolVar(&help, "parameters", false, "Display parameters information")
	flag.Parse()
	if help {
		if err := app.ParameterParser.Usage(app.Prefix, parameters...); err != nil {
			return err
		}
		return nil
	}
	if err := app.ParameterParser.Parse(app.Prefix, parameters...); err != nil {
		return err
	}
	for _, parameter := range parameters {
		if err := container.ProvideValue(parameter); err != nil {
			return fmt.Errorf("provide parameter failed; %w", err)
		}
	}
	// build bundle dependencies
	if err := buildBundles(container, sorted...); err != nil {
		return err
	}
	// resolve Logger from container
	// if Logger not found it will remain std
	if err = container.Resolve(&app.Logger); errors.Is(err, di.ErrTypeNotExists) {
		if err := container.ProvideValue(app.Logger, di.As(new(Logger))); err != nil {
			return err
		}
	}
	// start goroutine with os signal catch
	go app.catchSignals()
	startCtx, startCancel := context.WithTimeout(ctx, app.StartTimeout)
	// boot bundles
	hooks, err := beforeStart(startCtx, container, sorted...)
	startCancel()
	// if boot failed shutdown booted bundles
	if err != nil {
		// create context for shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), app.StopTimeout)
		defer cancel()
		if rserr := beforeShutdown(shutdownCtx, container, hooks); rserr != nil {
			return fmt.Errorf("%w (%s)", err, rserr)
		}
		printStartError(err)
	}
	app.Logger.Printf("Start")
	// run application, ignore context cancel error
	// default context lifecycle used for application shutdown
	if err := run(ctx, container); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	app.Logger.Printf("Stop")
	// create context for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), app.StopTimeout)
	defer cancel()
	// shutdown bundles in reverse order
	if err = beforeShutdown(shutdownCtx, container, hooks); err != nil {
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
	app.Logger.Printf(strings.Title(sign.String()))
	app.stop()
}
