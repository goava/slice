package testcmp

import (
	"context"
)

type FuncDispatcher struct {
	RunFunc func(ctx context.Context) error
}

func (d FuncDispatcher) Run(ctx context.Context) (err error) {
	return d.RunFunc(ctx)
}
