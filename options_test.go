package slice

import (
	"net/http"
	"testing"

	"github.com/goava/di"
	"github.com/stretchr/testify/require"
)

func TestRegisterBundles(t *testing.T) {
	s := New(
		Bundles(
			TestBundle{},
		),
	)
	require.Len(t, s.bundles, 1)
}

func TestDependencyInjection(t *testing.T) {
	s := New(
		DependencyInjection(
			di.Provide(func() *http.Server { return &http.Server{} }),
		),
	)
	require.Len(t, s.di, 1)
}
