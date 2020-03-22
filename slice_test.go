package slice

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/goava/di"
	"github.com/stretchr/testify/require"
)

func init() {
	stdLog = errStack{}
	exitError = func(err error) {
		panic(err.Error())
	}
}

func TestSlice_Start(t *testing.T) {
	t.Run("lifecycle", func(t *testing.T) {
		kernel := &KernelMock{
			RunFunc: func(ctx context.Context) error {
				return nil
			},
		}
		bundle := BootShutdownMock{
			DependencyInjectionFunc: func(builder ContainerBuilder) {
				builder.Provide(func(handler http.Handler) *http.Server {
					return &http.Server{
						Handler: handler,
					}
				})
				builder.Provide(func() *KernelMock {
					return kernel
				}, di.As(IKernel))
			},
			BootFunc: func(ctx context.Context, container Container) error {
				return nil
			},
			ShutdownFunc: func(ctx context.Context, container Container) error {
				return nil
			},
		}

		s := newLifecycle(
			Bundles(
				&bundle,
			),
			DependencyInjection(
				di.Provide(http.NewServeMux, di.As(new(http.Handler))),
			),
		)

		require.Nil(t, s.container)
		var mux *http.ServeMux
		require.NoError(t, s.initialization())
		require.NotNil(t, s.container)
		require.True(t, s.container.Has(&mux))

		var server *http.Server
		require.False(t, s.container.Has(&server))
		require.NoError(t, s.bundling())
		require.True(t, s.container.Has(&server))
		require.True(t, s.container.Has(IKernel))

		require.NoError(t, s.compiling()) // todo: check container compilation

		require.Nil(t, s.logger)
		require.Nil(t, s.kernel)
		require.NoError(t, s.resolving())
		require.NotNil(t, s.logger)
		require.NotNil(t, s.kernel)

		require.Len(t, bundle.BootCalls(), 0)
		s.boot()
		require.Len(t, bundle.BootCalls(), 1)

		require.Len(t, kernel.RunCalls(), 0)
		s.run()
		require.Len(t, kernel.RunCalls(), 1)

		require.Len(t, bundle.ShutdownCalls(), 0)
		s.shutdown()
		require.Len(t, bundle.ShutdownCalls(), 1)
	})
	t.Run("undefined kernel causes error", func(t *testing.T) {
		require.EqualError(t, newLifecycle().Start(), "slice.Kernel: not exists in container")
	})
	t.Run("failed initialization", func(t *testing.T) {
		s := newLifecycle(
			DependencyInjection(
				di.Provide(nil),
			),
		)
		err := s.Start()
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "di.Provide(..) failed:")
		require.Contains(t, err.Error(), "slice_test.go:93: constructor must be a function like func([dep1, dep2, ...]) (<result>, [cleanup, error]), got nil")
	})
	t.Run("failed bundling", func(t *testing.T) {
		bundle := &BundleMock{
			DependencyInjectionFunc: func(builder ContainerBuilder) {
				builder.Provide(nil)
				builder.Provide(nil)
			},
		}
		s := newLifecycle(
			Bundles(bundle),
		)
		require.EqualError(t, s.Start(), "*slice.BundleMock: Provide bundle components failed")
	})
}

type errStack struct {
	stack []string
}

func (s errStack) Error(err error) {
	s.stack = append(s.stack, fmt.Sprintf("error: %s", err.Error()))
}

func (s errStack) Fatal(err error) {
	panic(fmt.Sprintf("fatal: %s", err.Error()))
}

func (s errStack) Pop() string {
	defer func() {
		s.stack = s.stack[:len(s.stack)-1]
	}()
	return s.stack[len(s.stack)-1]
}
