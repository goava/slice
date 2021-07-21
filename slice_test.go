package slice_test

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/goava/di"
	"github.com/goava/slice/bundle"
	"github.com/goava/slice/testcmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/goava/slice"
)

func TestInitializationErrors(t *testing.T) {
	t.Run("application name must be specified rerun", func(t *testing.T) {
		if os.Getenv("APP_TEST_CRASH") == "1" {
			slice.Run()
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestInitializationErrors")
		cmd.Env = append(os.Environ(), "APP_TEST_CRASH=1")
		_, err := cmd.Output()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			require.EqualError(t, e, "exit status 1")
			require.Contains(t, string(e.Stderr), "application name must be specified, see slice.SetName() option")
			return
		}
		t.Fatalf("process started with err %v, want exit status 1", err)
	})

	t.Run("application name must be specified", func(t *testing.T) {
		logger := &testcmp.Log{}
		require.PanicsWithValue(t, "fatal interruption", func() {
			slice.Run(
				slice.WithLogger(logger),
			)
		}, "app should stop with panic")
		require.Len(t, logger.FatalLogs, 1, "logger should have 1 fatal message")
		require.Equal(t, "application name must be specified, see slice.SetName() option", logger.FatalLogs[0])
	})

	t.Run("bundle without name cause error", func(t *testing.T) {
		logger := &testcmp.FmtLog{}
		slice.Run(
			slice.WithName("app"),
			slice.WithLogger(logger),
			slice.WithBundles(bundle.New()),
		)
		require.Len(t, logger.FatalLogs, 1)
		require.Equal(t, "prepare bundles: bundle with index 0: empty name", logger.FatalLogs[0])
	})

	t.Run("invalid component causes error", func(t *testing.T) {
		logger := &testcmp.FmtLog{}
		slice.Run(
			slice.WithName("app"),
			slice.WithLogger(logger),
			slice.WithComponents(
				slice.Provide(nil),
			),
		)
		require.Len(t, logger.FatalLogs, 1)
		require.Contains(t, logger.FatalLogs[0], "initialization: create container failed: ")
		require.Contains(t, logger.FatalLogs[0], ": invalid constructor signature, got nil")
	})
}

func TestDefaultComponents(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = []string{"app"}
	_ = os.Setenv("ENV", "")
	_ = os.Setenv("DEBUG", "")
	t.Run("default components provided", func(t *testing.T) {
		called := false
		dispatcher := func(ctx context.Context, info slice.Info, env slice.Env) *testcmp.FuncDispatcher {
			return &testcmp.FuncDispatcher{RunFunc: func(ctx context.Context) error {
				require.NotNil(t, ctx)
				require.NotNil(t, info)
				require.Equal(t, "prod", info.Env.String())
				require.False(t, info.Debug)
				require.Equal(t, "app", info.Name)
				require.Equal(t, "prod", env.String())
				called = true
				return nil
			}}
		}
		slice.Run(
			slice.WithName("app"),
			slice.WithComponents(
				slice.Provide(dispatcher, di.As(new(slice.Dispatcher))),
			),
		)
		require.True(t, called)
	})
}

func TestProvideBundleComponents(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = []string{"app"}
	_ = os.Setenv("ENV", "")
	_ = os.Setenv("DEBUG", "")
	t.Run("bundle components provided on start", func(t *testing.T) {
		called := false
		dispatcher := func(s string, i int32) *testcmp.FuncDispatcher {
			return &testcmp.FuncDispatcher{RunFunc: func(ctx context.Context) error {
				assert.Equal(t, "test-string", s)
				assert.Equal(t, int32(1), i)
				called = true
				return nil
			}}
		}
		stringBundle := bundle.New(
			bundle.WithName("string-bundle"),
			bundle.WithComponents(
				slice.Supply("test-string"),
			),
		)
		int32Bundle := bundle.New(
			bundle.WithName("int32-bundle"),
			bundle.WithComponents(
				slice.Supply(int32(1)),
			),
		)
		slice.Run(
			slice.WithName("app"),
			slice.WithBundles(
				stringBundle,
				int32Bundle,
			),
			slice.WithComponents(
				slice.Provide(dispatcher, di.As(new(slice.Dispatcher))),
			),
		)
		require.True(t, called)
	})
}

type TestDispatcher struct {
}

func NewTestDispatcher() *TestDispatcher {
	return &TestDispatcher{}
}

func (r TestDispatcher) Run(ctx context.Context) error {
	return nil
}

func (r TestDispatcher) Stop() error {
	return nil
}
