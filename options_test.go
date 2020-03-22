package slice

import (
	"context"
	"net/http"
	"testing"

	"github.com/goava/di"
	"github.com/stretchr/testify/require"
)

func TestRegisterBundles(t *testing.T) {
	bundle := &BootShutdownMock{
		DependencyInjectionFunc: func(builder ContainerBuilder) {},
		BootFunc:                func(ctx context.Context, container Container) error { return nil },
		ShutdownFunc:            func(ctx context.Context, container Container) error { return nil },
	}

	s := newLifecycle(
		Bundles(bundle),
	)
	require.Len(t, s.bundles, 1)
}

func TestDependencyInjection(t *testing.T) {
	s := newLifecycle(
		DependencyInjection(
			di.Provide(func() *http.Server { return &http.Server{} }),
		),
	)
	require.Len(t, s.di, 1)
}
