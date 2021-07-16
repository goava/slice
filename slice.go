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

type state int

const (
	none state = iota
	initialization
	configuring
	starting
	running
	shutdown
)

// Application is a control part of application.
type Application struct {
	Name            string
	Prefix          string
	Parameters      []Parameter
	Dispatcher      Dispatcher
	Bundles         []Bundle
	StartTimeout    time.Duration
	ShutdownTimeout time.Duration
	Logger          Logger
	ParameterParser ParameterParser

	// providers contains type providers. Only slice.Provide() and slice.Supply() works.
	providers []di.Option
	env       Env
	debug     bool
	state     state
	stop      func()
}

// Start starts application.
func (app *Application) Start() error {
	// STATE: INITIALIZATION
	app.state = initialization
	if app.Logger == nil {
		app.Logger = &stdLogger{} // std logger logs messages before container initialization
	}
	if len(app.Name) == 0 {
		return fmt.Errorf("application name must be specified, see slice.SetName() option")
	}
	// initialize context
	base, stop := context.WithCancel(context.Background())
	app.stop = stop
	ctx := NewContext(base)
	// lookup environment
	env, _ := lookupEnv(defaultEnv)
	app.env = parseEnv(env)
	if debug, ok := lookupEnv(defaultDebug); ok {
		app.debug = strings.ToLower(debug) == "true"
	}
	// build app info
	info := Info{
		Name:  app.Name,
		Env:   app.env,
		Debug: app.debug,
	}
	// check bundle acyclic and sort dependencies
	sorted, ok := prepareBundles(app.Bundles)
	if !ok {
		return fmt.Errorf("bundle cyclic detected") // todo: improve error message
	}
	// prepare bundle components
	for _, bundle := range sorted {
		bundle.apply(app)
	}
	// prepare application components
	providers := []di.Option{
		di.Provide(func() *Context { return ctx }, di.As(new(context.Context))),
		di.Provide(func() Info { return info }),
	}
	providers = append(providers, app.providers...)
	// validate container with all application components
	container, err := createContainer(providers...)
	if err != nil {
		return fmt.Errorf("initialization: %w", err)
	}
	// STATE: CONFIGURING
	app.state = configuring
	if app.ParameterParser == nil {
		err = container.Resolve(&app.ParameterParser)
		if err != nil && !errors.Is(err, di.ErrTypeNotExists) {
			return fmt.Errorf("configuring: parameter parser: %w", err)
		}
		if err != nil && errors.Is(err, di.ErrTypeNotExists) {
			app.ParameterParser = &stdParameterParser{}
		}
	}
	// collect application parameters
	parameters := app.Parameters
	// add bundle parameters
	for _, bundle := range sorted {
		parameters = append(parameters, bundle.Parameters...)
	}
	// check parameters
	var parametersFlag bool
	flag.BoolVar(&parametersFlag, "parameters", false, "Display parameters information")
	flag.Parse()
	if parametersFlag {
		if err := app.ParameterParser.Usage(app.Prefix, parameters...); err != nil {
			return fmt.Errorf("configuring: usage: %w", err)
		}
		return nil
	}
	// parameter parser decorator, implemented for lazy parameter loading
	parseParameters := func(pointer di.Value) error {
		if err := app.ParameterParser.Parse(app.Prefix, pointer); err != nil {
			return fmt.Errorf("configuring: parse: %w", err)
		}
		return nil
	}
	for _, parameter := range parameters {
		if err := container.ProvideValue(parameter, di.Decorate(parseParameters)); err != nil {
			return fmt.Errorf("configuring: parameters: %w", err)
		}
	}
	// resolve logger
	err = container.Resolve(&app.Logger)
	if err != nil && errors.Is(err, di.ErrTypeNotExists) {
		if err := container.ProvideValue(app.Logger, di.As(new(Logger))); err != nil {
			return fmt.Errorf("configuring: logger: %w", err)
		}
	}
	if err != nil && !errors.Is(err, di.ErrTypeNotExists) {
		return fmt.Errorf("configuring: logger: %w", err)
	}
	// STATE: STARTING
	app.state = starting
	var dispatchers []Dispatcher
	has, err := container.Has(&dispatchers)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("no one slice.Dispatcher found")
	}
	// set timeouts
	if app.StartTimeout == 0 {
		app.StartTimeout = defaultTimeout
	}
	if app.ShutdownTimeout == 0 {
		app.ShutdownTimeout = defaultTimeout
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
		shutdownCtx, cancel := context.WithTimeout(context.Background(), app.ShutdownTimeout)
		defer cancel()
		if rserr := beforeShutdown(shutdownCtx, container, hooks); rserr != nil {
			return fmt.Errorf("%w (%s)", err, rserr)
		}
		printStartError(err)
	}
	app.Logger.Printf("Start")
	// resolve dispatchers
	if err := container.Resolve(&dispatchers); err != nil {
		return fmt.Errorf("dispatch failed: %w", err)
	}
	// STATE: RUNNING
	app.state = running
	// dispatch application, ignore context cancel error
	// default context lifecycle used for application shutdown
	if err := dispatch(ctx, stop, dispatchers); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	// STATE: SHUTDOWN
	app.state = shutdown
	app.Logger.Printf("Stop")
	// create context for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), app.ShutdownTimeout)
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
