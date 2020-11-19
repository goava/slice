package slice_test

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/goava/di"
	"github.com/goava/slice"
	"github.com/stretchr/testify/require"
)

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

func TestInitializationErrors(t *testing.T) {
	t.Run("application name must be specified", func(t *testing.T) {
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
		t.Fatalf("process ran with err %v, want exit status 1", err)
	})
}

func TestRun(t *testing.T) {
	t.Run("run", func(t *testing.T) {
		slice.Run(
			slice.SetName("test_run"),
			slice.RegisterBundles(),
			slice.ConfigureContainer(
				di.Provide(NewTestDispatcher, di.As(new(slice.Dispatcher))),
			),
		)
	})
}
