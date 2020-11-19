package slice

import (
	"net/http"
	"testing"

	"github.com/goava/di"
	"github.com/stretchr/testify/require"
)

func TestBundleContainerBuilder_Has(t *testing.T) {
	c, err := di.New(di.Provide(http.NewServeMux))
	require.NoError(t, err)
	require.NotNil(t, c)
	cb := containerBuilder{container: c}
	var mux *http.ServeMux
	has, err := cb.Has(&mux)
	require.NoError(t, err)
	require.True(t, has)
	var server *http.Server
	has, err = cb.Has(&server)
	require.NoError(t, err)
	require.False(t, has)
}

func TestBundleContainerBuilder_Provide(t *testing.T) {
	t.Run("builder provide component to container", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		cb := containerBuilder{container: c}
		cb.Provide(func() *http.Server { return &http.Server{} })
		require.Len(t, cb.errs, 0)
		var server *http.Server
		has, err := cb.Has(&server)
		require.NoError(t, err)
		require.True(t, has)
	})

	t.Run("if provide error builder saves error", func(t *testing.T) {
		c, err := di.New()
		require.NoError(t, err)
		require.NotNil(t, c)
		cb := containerBuilder{container: c}
		cb.Provide(func() {})
		require.Len(t, cb.errs, 1)
		require.Error(t, cb.Error())
		require.Contains(t, cb.Error().Error(), "invalid constructor signature, got func()")
	})
}

func TestInspectBundles(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		bundles := []Bundle{
			FirstBundle{},
			SecondBundle{},
		}
		inspected := inspectBundles(bundles...)
		for _, entry := range inspected {
			require.Equal(t, entry.name, bundleName(entry.Bundle))
		}
	})
}
