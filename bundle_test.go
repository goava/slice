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
	require.True(t, cb.Has(&mux))
	var server *http.Server
	require.False(t, cb.Has(&server))
}

func TestBundleContainerBuilder_Provide(t *testing.T) {
	t.Run("builder provide component to container", func(t *testing.T) {
		c, err := di.New(
			di.WithCompile(),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		cb := containerBuilder{container: c}
		cb.Provide(func() *http.Server { return &http.Server{} })
		require.Len(t, cb.errs, 0)
		var server *http.Server
		require.True(t, cb.Has(&server))
	})

	t.Run("if provide error builder saves error", func(t *testing.T) {
		c, err := di.New(
			di.WithCompile(),
		)
		require.NoError(t, err)
		require.NotNil(t, c)
		cb := containerBuilder{container: c}
		cb.Provide(func() {})
		require.Len(t, cb.errs, 1)
		require.Error(t, cb.Error())
		require.Contains(t, cb.Error().Error(), "invalid constructor signature, got func()")
	})
}
