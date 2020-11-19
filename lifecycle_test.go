package slice

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/goava/di"
	"github.com/stretchr/testify/require"
)

func TestLifecycle_initialization(t *testing.T) {
	t.Run("provide user dependency", func(t *testing.T) {
		c, err := createContainer(
			di.Provide(http.NewServeMux),
		)
		require.NoError(t, err)
		var mux *http.ServeMux
		require.True(t, c.Has(&mux))
	})

	t.Run("incorrect option cause error", func(t *testing.T) {
		c, err := createContainer(
			di.Provide(func() {}),
		)
		require.Nil(t, c)
		require.Error(t, err)
		require.Contains(t, err.Error(), "lifecycle_test.go:")
		require.Contains(t, err.Error(), ": invalid constructor signature, got func()")
	})
}

func TestLifecycle_configureBundles(t *testing.T) {
	t.Run("process iterates over all bundles", func(t *testing.T) {
		bundles := []bundle{
			{
				name:   "first-bundle",
				Bundle: &BundleMock{},
			},
			{
				name:   "second-bundle",
				Bundle: &BundleMock{},
			},
			{
				name:   "third-bundle",
				Bundle: &BundleMock{},
			},
		}
		i := 0
		err := configureBundles(func(bundle Bundle) error {
			require.Equal(t, bundles[i].Bundle, bundle)
			i++
			return nil
		}, bundles...)
		require.NoError(t, err)
		require.Equal(t, 3, i)
	})

	t.Run("process error causes configure error", func(t *testing.T) {
		bundle := bundle{
			name:   "error-bundle",
			Bundle: &BundleMock{},
		}
		err := configureBundles(func(bundle Bundle) error {
			return errors.New("unexpected error")
		}, bundle)
		require.EqualError(t, err, "configure error-bundle bundle failed: unexpected error")
	})
}

func TestLifecycle_buildBundles(t *testing.T) {
	t.Run("bundle builds in correct order", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		var order []string
		firstBundle := &BundleMock{
			BuildFunc: func(builder ContainerBuilder) {
				order = append(order, "first")
			},
		}
		secondBundle := &BundleMock{
			BuildFunc: func(builder ContainerBuilder) {
				order = append(order, "second")
			},
		}
		bundles := []bundle{
			{
				name:   "first-bundle",
				Bundle: firstBundle,
			},
			{
				name:   "second-bundle",
				Bundle: secondBundle,
			},
		}
		err = buildBundles(c, bundles...)
		require.NoError(t, err)
		require.Len(t, firstBundle.BuildCalls(), 1)
		require.Len(t, secondBundle.BuildCalls(), 1)
		require.Equal(t, []string{"first", "second"}, order)
	})

	t.Run("bundle build error return as one", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		errorBundle := bundle{
			name: "error-bundle",
			Bundle: &BundleMock{
				BuildFunc: func(builder ContainerBuilder) {
					builder.Provide(func() {})
					builder.Provide(nil)
					builder.Provide(struct{}{})
				},
			},
		}
		err = buildBundles(c, errorBundle)
		require.Error(t, err)
		require.Contains(t, err.Error(), "build error-bundle bundle failed:")
		require.Contains(t, err.Error(), "invalid constructor signature, got func()")
		require.Contains(t, err.Error(), "invalid constructor signature, got nil")
		require.Contains(t, err.Error(), "invalid constructor signature, got struct {}")
	})
}

func TestLifecycle_boot(t *testing.T) {
	t.Run("iterates over bundles and run boot function", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		var order []string
		firstBundle := &BootShutdownMock{
			BootFunc: func(ctx context.Context, container Container) error {
				order = append(order, "first-bundle")
				return nil
			},
		}
		secondBundle := &BootShutdownMock{
			BootFunc: func(ctx context.Context, container Container) error {
				order = append(order, "second-bundle")
				return nil
			},
		}
		bundles := []bundle{
			{
				name:   "first-bundle",
				Bundle: firstBundle,
			},
			{
				name:   "second-bundle",
				Bundle: secondBundle,
			},
		}
		shutdowns, err := boot(context.Background(), c, bundles...)
		require.NoError(t, err)
		require.Len(t, shutdowns, 2)
		require.Len(t, firstBundle.BootCalls(), 1)
		require.Len(t, secondBundle.BootCalls(), 1)
		require.Equal(t, []string{"first-bundle", "second-bundle"}, order)
	})

	t.Run("bundle boot error causes boot error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		bundle := bundle{
			name: "error-bundle",
			Bundle: &BootShutdownMock{
				BootFunc: func(ctx context.Context, container Container) error {
					return errors.New("unexpected error")
				},
			},
		}
		shutdowns, err := boot(context.Background(), c, bundle)
		require.EqualError(t, err, "boot error-bundle bundle failed: unexpected error")
		require.Len(t, shutdowns, 0)
	})

	t.Run("shutdowns correct on context cancel", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		firstBundle := &BootShutdownMock{
			BootFunc: func(ctx context.Context, container Container) error {
				time.Sleep(2 * time.Millisecond)
				return nil
			},
		}
		secondBundle := &BootShutdownMock{
			BootFunc: func(ctx context.Context, container Container) error {
				return nil
			},
		}
		bundles := []bundle{
			{
				name:   "first-bundle",
				Bundle: firstBundle,
			},
			{
				name:   "second-bundle",
				Bundle: secondBundle,
			},
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		shutdowns, err := boot(ctx, c, bundles...)
		require.EqualError(t, err, "boot first-bundle bundle failed: context canceled")
		require.Len(t, firstBundle.BootCalls(), 0)
		require.Len(t, secondBundle.BootCalls(), 0)
		require.Len(t, shutdowns, 0)
	})
}

func TestLifecycle_drun(t *testing.T) {
	t.Run("resolve dispatcher and run", func(t *testing.T) {
		dispatcher := &DispatcherMock{
			RunFunc: func(ctx context.Context) error {
				return nil
			},
		}
		c, err := di.New(di.Provide(func() *DispatcherMock { return dispatcher }, di.As(new(Dispatcher))))
		require.NotNil(t, c)
		err = run(context.Background(), c)
		require.NoError(t, err)
		require.Len(t, dispatcher.RunCalls(), 1)
	})

	t.Run("undefined dispatcher cause error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		err = run(context.Background(), c)
		require.Error(t, err)
		require.Contains(t, err.Error(), "resolve dispatcher failed: ")
		require.Contains(t, err.Error(), "lifecycle.go:")
		require.Contains(t, err.Error(), ": type slice.Dispatcher not exists in the container")
	})

	t.Run("run error causes error", func(t *testing.T) {
		dispatcher := &DispatcherMock{
			RunFunc: func(ctx context.Context) error {
				return errors.New("unexpected error")
			},
		}
		c, err := di.New(di.Provide(func() *DispatcherMock { return dispatcher }, di.As(new(Dispatcher))))
		require.NotNil(t, c)
		err = run(context.Background(), c)
		require.EqualError(t, err, "failure: unexpected error")
		require.Len(t, dispatcher.RunCalls(), 1)
	})
}

func TestLifecycle_reverseShutdown(t *testing.T) {
	t.Run("reverse order", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		var order []string
		shutdowns := shutdowns{
			{
				name: "first-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					order = append(order, "first-shutdown")
					return nil
				},
			},
			{
				name: "second-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					order = append(order, "second-shutdown")
					return nil
				},
			},
			{
				name: "third-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					order = append(order, "third-shutdown")
					return nil
				},
			},
		}
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		err = reverseShutdown(ctx, c, shutdowns)
		require.NoError(t, err)
		require.Equal(t, []string{"third-shutdown", "second-shutdown", "first-shutdown"}, order)
	})

	t.Run("shutdown errors returns one", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		shutdowns := shutdowns{
			{
				name: "first-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					return errors.New("first-error")
				},
			},
			{
				name: "second-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					return errors.New("second-error")
				},
			},
			{
				name: "third-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					return errors.New("third-error")
				},
			},
		}
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		err = reverseShutdown(ctx, c, shutdowns)
		require.EqualError(t, err, "shutdown failed: shutdown third-shutdown failed: third-error; shutdown second-shutdown failed: second-error; shutdown first-shutdown failed: first-error")
	})

	t.Run("context error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		shutdowns := shutdowns{
			{
				name: "first-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					time.Sleep(time.Millisecond)
					return nil
				},
			},
			{
				name: "second-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					return errors.New("second-error")
				},
			},
			{
				name: "third-shutdown",
				shutdown: func(ctx context.Context, container Container) error {
					return errors.New("third-error")
				},
			},
		}
		ctx, _ := context.WithTimeout(context.Background(), time.Nanosecond)
		err = reverseShutdown(ctx, c, shutdowns)
		require.EqualError(t, err, "shutdown failed: context deadline exceeded")
	})
}
