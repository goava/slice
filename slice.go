package slice

import (
	"context"

	"github.com/goava/di"
)

var (
	// IKernel is a argument for di.As() container option.
	IKernel = new(Kernel)
	ILogger = new(Logger)
)

// Run runs slice application.
func New(options ...Option) *Slice {
	s := Slice{}
	for _, opt := range options {
		opt.apply(&s)
	}
	return &s
}

// Slice is a component-based framework. It built around dependency injection.
type Slice struct {
	di        []di.Option
	bundles   []Bundle
	booted    []BootShutdown
	container *di.Container
	logger    Logger
	kernel    Kernel
}

// Starts start slice.
func (s *Slice) Start() {
	s.initialization()
	s.bundling()
	s.compiling()
	s.resolving()
	s.boot()
	s.run()
	s.shutdown()
}

func (s *Slice) initialization() {
	s.di = append(s.di, di.WithCompile())
	container, err := di.New(s.di...)
	if err != nil {
		exitError(err)
	}
	s.container = container
}

func (s *Slice) bundling() {
	for _, b := range s.bundles {
		builder := &bundleContainerBuilder{
			container: s.container,
			bundleErr: bundleDIErrors{bundle: b},
		}
		b.DependencyInjection(builder)
		if len(builder.bundleErr.list) > 0 {
			exitError(builder.bundleErr)
		}
	}
}

func (s *Slice) compiling() {
	if err := s.container.Compile(); err != nil {
		exitError(err)
	}
}

func (s *Slice) resolving() {
	if err := s.container.Resolve(&s.logger); err != nil {
		s.logger = stdLog
	}
	if err := s.container.Resolve(&s.kernel); err != nil {
		exitError(err)
	}
}

func (s *Slice) boot() {
	err := errBootFailed{}
	for _, b := range s.bundles {
		if bs, ok := b.(BootShutdown); ok {
			if bootErr := bs.Boot(context.TODO(), s.container); bootErr != nil {
				err = append(err, bootErr)
				continue
			}
			s.booted = append(s.booted, bs)
		}
	}
	if len(err) > 0 {
		s.shutdown()
		s.logger.Fatal(err)
	}
}

func (s *Slice) run() {
	ctx, _ := context.WithCancel(context.Background())
	if err := s.kernel.Run(ctx); err != nil {
		s.logger.Error(err)
	}
}

func (s *Slice) shutdown() {
	err := errShutdownFailed{}
	for _, shutdown := range s.booted {
		if shutdownErr := shutdown.Shutdown(context.TODO(), s.container); shutdownErr != nil {
			err = append(err, shutdownErr)
		}
	}
	if len(err) > 0 {
		s.logger.Error(err)
	}
}
