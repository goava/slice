package slice_test

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/goava/di"
	"github.com/goava/slice/testcmp"
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
}

func TestRun(t *testing.T) {
	t.Run("full example", func(t *testing.T) {
		slice.Run(
			slice.WithName("test_run"),
			slice.WithComponents(
				slice.Provide(NewTestDispatcher, di.As(new(slice.Dispatcher))),
			),
		)
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
