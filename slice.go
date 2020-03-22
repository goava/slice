package slice

import (
	"context"
	"fmt"

	"github.com/goava/di"
)

var (
	// IKernel is a argument for di.As() container option.
	IKernel = new(Kernel)
	ILogger = new(Logger)
)

// Run runs slice lifecycle with options.
func Run(options ...Option) {
	lf := newLifecycle(options...)
	if err := lf.Start(); err != nil {
		exitError(err)
	}
}

// newLifecycle creates slice application.
func newLifecycle(options ...Option) *lifecycle {
	s := lifecycle{}
	for _, opt := range options {
		opt.apply(&s)
	}
	return &s
}

// lifecycle is a component-based framework. It built around dependency injection.
type lifecycle struct {
	di        []di.Option
	bundles   []Bundle
	booted    []BootShutdown
	container *di.Container
	logger    Logger
	kernel    Kernel
}

// Starts start slice.
func (s *lifecycle) Start() error {
	if err := s.initialization(); err != nil {
		return err
	}
	if err := s.bundling(); err != nil {
		return err
	}
	if err := s.compiling(); err != nil {
		return err
	}
	if err := s.resolving(); err != nil {
		return err
	}
	if err := s.boot(); err != nil {
		s.shutdown()
		return err
	}
	s.run()
	s.shutdown()
	return nil
}

func (s *lifecycle) initialization() error {
	s.di = append(s.di, di.WithCompile())
	container, err := di.New(s.di...)
	if err != nil {
		return err
	}
	s.container = container
	return nil
}

func (s *lifecycle) bundling() error {
	for _, b := range s.bundles {
		builder := &bundleContainerBuilder{
			container: s.container,
			bundleErr: bundleDIErrors{bundle: b},
		}
		b.DependencyInjection(builder)
		if len(builder.bundleErr.list) > 0 {
			return builder.bundleErr
		}
	}
	return nil
}

func (s *lifecycle) compiling() error {
	if err := s.container.Compile(); err != nil {
		return err
	}
	return nil
}

func (s *lifecycle) resolving() error {
	if err := s.container.Resolve(&s.logger); err != nil {
		s.logger = stdLog
	}
	if err := s.container.Resolve(&s.kernel); err != nil {
		return err
	}
	return nil
}

func (s *lifecycle) boot() error {
	for _, b := range s.bundles {
		if bs, ok := b.(BootShutdown); ok {
			if err := bs.Boot(context.TODO(), s.container); err != nil {
				return fmt.Errorf("boot failed: %s", err)
			}
			s.booted = append(s.booted, bs)
		}
	}
	return nil
}

func (s *lifecycle) run() {
	ctx, _ := context.WithCancel(context.Background())
	if err := s.kernel.Run(ctx); err != nil {
		s.logger.Error(err)
	}
}

func (s *lifecycle) shutdown() {
	for i := len(s.booted) - 1; i >= 0; i-- {
		shutdown := s.booted[i].Shutdown
		if err := shutdown(context.TODO(), s.container); err != nil {
			s.logger.Error(err)
		}
	}
}
