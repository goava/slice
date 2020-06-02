package slice

import (
	"context"
	"testing"

	"github.com/goava/di"
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

func TestRun(t *testing.T) {
	t.Run("run", func(t *testing.T) {
		Run(
			RegisterBundles(),
			ConfigureContainer(
				di.Provide(NewTestDispatcher, di.As(new(Dispatcher))),
			),
		)
	})
}
