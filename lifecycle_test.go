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

func TestLifecycle_createContainer(t *testing.T) {
	t.Run("provide user dependency", func(t *testing.T) {
		c, err := createContainer(
			di.Provide(http.NewServeMux),
		)
		require.NoError(t, err)
		var mux *http.ServeMux
		has, err := c.Has(&mux)
		require.NoError(t, err)
		require.True(t, has)
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

func TestLifecycle_buildBundles(t *testing.T) {
	t.Run("bundle components provided in correct order", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		var order []string
		first := Bundle{
			Name: "first-bundle",
			Components: []di.Option{
				di.Invoke(func() {
					order = append(order, "first")
				}),
			},
		}
		second := Bundle{
			Name: "second-bundle",
			Components: []di.Option{
				di.Invoke(func() {
					order = append(order, "second")
				}),
			},
		}
		err = buildBundles(c, first, second)
		require.NoError(t, err)
		require.Equal(t, []string{"first", "second"}, order)
	})

	t.Run("bundle build error return as one", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		errorBundle := Bundle{
			Name: "error-bundle",
			Components: []di.Option{
				di.Provide(func() {}),
				di.Provide(nil),
				di.Provide(struct{}{}),
			},
		}
		err = buildBundles(c, errorBundle)
		require.Error(t, err)
		require.Contains(t, err.Error(), "build error-bundle bundle failed:")
		require.Contains(t, err.Error(), "invalid constructor signature, got func()")
		// require.Contains(t, err.Error(), "invalid constructor signature, got nil")
		// require.Contains(t, err.Error(), "invalid constructor signature, got struct {}")
	})
}

func TestLifecycle_before(t *testing.T) {
	t.Run("iterates over bundles and run before hook", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		var order []string
		firstBundle := Bundle{
			Name: "first-bundle",
			Hooks: []Hook{{
				Before: func() {
					order = append(order, "first-bundle")
				},
				After: func() {},
			}},
		}
		secondBundle := Bundle{
			Name: "second-bundle",
			Hooks: []Hook{{
				Before: func() {
					order = append(order, "second-bundle")
				},
			}},
		}
		shutdowns, err := beforeStart(context.Background(), c, firstBundle, secondBundle)
		require.NoError(t, err)
		require.Len(t, shutdowns, 1)
		require.Equal(t, []string{"first-bundle", "second-bundle"}, order)
	})

	t.Run("bundle boot error causes boot error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		bundle := Bundle{
			Name: "error-bundle",
			Hooks: []Hook{{
				Before: func() error { return errors.New("unexpected error") },
			}},
		}
		hooks, err := beforeStart(context.Background(), c, bundle)
		require.EqualError(t, err, "- boot error-bundle bundle failed: unexpected error\n")
		require.Len(t, hooks, 0)
	})

	t.Run("shutdowns correct on context cancel", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		firstBundle := Bundle{
			Name: "first-bundle",
			Hooks: []Hook{{
				Before: func() error {
					time.Sleep(2 * time.Millisecond)
					return nil
				},
			}},
		}
		secondBundle := Bundle{
			Name: "second-bundle",
			Hooks: []Hook{{
				Before: func() {},
			}},
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		hooks, err := beforeStart(ctx, c, firstBundle, secondBundle)
		require.EqualError(t, err, "boot first-bundle bundle failed: context canceled")
		require.Len(t, hooks, 0)
	})
}

func TestLifecycle_dispatch(t *testing.T) {
	t.Run("resolve dispatchers and run only once", func(t *testing.T) {
		dispatcher := &DispatcherMock{
			RunFunc: func(ctx context.Context) error {
				return nil
			},
		}
		ctx, cancel := context.WithCancel(context.Background())
		err := dispatch(ctx, cancel, []Dispatcher{dispatcher})
		require.NoError(t, err)
		require.Len(t, dispatcher.RunCalls(), 1)
	})

	// todo: rewrite test case
	//t.Run("undefined dispatcher cause error", func(t *testing.T) {
	//	c, err := di.New()
	//	require.NoError(t, err)
	//	require.NotNil(t, c)
	//	ctx, cancel := context.WithCancel(context.Background())
	//	err = dispatch(ctx, cancel)
	//	require.Error(t, err)
	//	require.Contains(t, err.Error(), "dispatch failed: ")
	//	require.Contains(t, err.Error(), "lifecycle.go:")
	//	require.Contains(t, err.Error(), ": type []slice.Dispatcher not exists in the container")
	//})

	t.Run("run error causes error", func(t *testing.T) {
		d1 := &DispatcherMock{
			RunFunc: func(ctx context.Context) error {
				return errors.New("unexpected error")
			},
		}
		contextCancelled := false
		d2 := &DispatcherMock{
			RunFunc: func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					contextCancelled = true
					return ctx.Err()
				case <-time.After(time.Second):
					require.Fail(t, "context should be cancelled")
					return nil
				}
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		err := dispatch(ctx, cancel, []Dispatcher{d1, d2})
		require.EqualError(t, err, "failure: unexpected error")
		require.Len(t, d1.RunCalls(), 1)
		require.True(t, contextCancelled)
	})
}

func TestLifecycle_after(t *testing.T) {
	t.Run("reverse order", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		var order []string
		hooks := []hook{
			{
				name: "first-shutdown",
				hook: func() {
					order = append(order, "first-shutdown")
				},
			},
			{
				name: "second-shutdown",
				hook: func() {
					order = append(order, "second-shutdown")
				},
			},
			{
				name: "third-shutdown",
				hook: func() {
					order = append(order, "third-shutdown")
				},
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err = beforeShutdown(ctx, c, hooks)
		require.NoError(t, err)
		require.Equal(t, []string{"third-shutdown", "second-shutdown", "first-shutdown"}, order)
	})

	t.Run("shutdown errors returns one", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		hooks := []hook{
			{
				name: "first-shutdown",
				hook: func() error {
					return errors.New("first-error")
				},
			},
			{
				name: "second-shutdown",
				hook: func() error {
					return errors.New("second-error")
				},
			},
			{
				name: "third-shutdown",
				hook: func() error {
					return errors.New("third-error")
				},
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err = beforeShutdown(ctx, c, hooks)
		require.EqualError(t, err, "shutdown failed: shutdown third-shutdown failed: third-error; shutdown second-shutdown failed: second-error; shutdown first-shutdown failed: first-error")
	})

	t.Run("context error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		shutdowns := []hook{
			{
				name: "first-shutdown",
				hook: func() error {
					time.Sleep(time.Hour)
					return nil
				},
			},
			{
				name: "second-shutdown",
				hook: func() error {
					return errors.New("second-error")
				},
			},
			{
				name: "third-shutdown",
				hook: func() error {
					return errors.New("third-error")
				},
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()
		err = beforeShutdown(ctx, c, shutdowns)
		require.EqualError(t, err, "shutdown failed: context deadline exceeded")
	})
}
