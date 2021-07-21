package slice

import (
	"context"
	"os"
	"testing"

	"github.com/goava/di"
	"github.com/goava/slice/testcmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnv_IsDev(t *testing.T) {
	tests := []struct {
		name string
		e    Env
		want bool
	}{
		{
			name: "dev mode",
			e:    "dev",
			want: true,
		},
		{
			name: "not dev mode",
			e:    "prod",
			want: false,
		},
		{
			name: "dev mode uppercase",
			e:    "DEV",
			want: true,
		},
		{
			name: "dev mode postfix",
			e:    "dev-feature-1",
			want: true,
		},
		{
			name: "dev mode postfix uppercase",
			e:    "DEV-FEATURE-1",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsDev(); got != tt.want {
				t.Errorf("IsDev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = []string{"app"}
	t.Run("default production environment", func(t *testing.T) {
		called := false
		dispatcher := func(env Env) *testcmp.FuncDispatcher {
			return &testcmp.FuncDispatcher{RunFunc: func(ctx context.Context) error {
				assert.True(t, env.String() == "prod")
				assert.False(t, env.IsTest())
				assert.False(t, env.IsDev())
				called = true
				return nil
			}}
		}
		Run(
			WithName("app"),
			WithComponents(
				Provide(dispatcher, di.As(new(Dispatcher))),
			),
		)
		require.True(t, called)
	})

	t.Run("development environment", func(t *testing.T) {
		_ = os.Setenv("ENV", "dev")
		called := false
		dispatcher := func(env Env) *testcmp.FuncDispatcher {
			return &testcmp.FuncDispatcher{RunFunc: func(ctx context.Context) error {
				assert.True(t, env.String() == "dev")
				assert.False(t, env.IsTest())
				assert.True(t, env.IsDev())
				called = true
				return nil
			}}
		}
		Run(
			WithName("app"),
			WithComponents(
				Provide(dispatcher, di.As(new(Dispatcher))),
			),
		)
		require.True(t, called)
	})

	t.Run("testing environment", func(t *testing.T) {
		_ = os.Setenv("ENV", "test")
		called := false
		dispatcher := func(env Env) *testcmp.FuncDispatcher {
			return &testcmp.FuncDispatcher{RunFunc: func(ctx context.Context) error {
				assert.True(t, env.String() == "test")
				assert.True(t, env.IsTest())
				assert.False(t, env.IsDev())
				called = true
				return nil
			}}
		}
		Run(
			WithName("app"),
			WithComponents(
				Provide(dispatcher, di.As(new(Dispatcher))),
			),
		)
		require.True(t, called)
	})

	t.Run("development prefix environment", func(t *testing.T) {
		_ = os.Setenv("ENV", "dev-1")
		called := false
		dispatcher := func(env Env) *testcmp.FuncDispatcher {
			return &testcmp.FuncDispatcher{RunFunc: func(ctx context.Context) error {
				assert.True(t, env.String() == "dev-1")
				assert.False(t, env.IsTest())
				assert.True(t, env.IsDev())
				called = true
				return nil
			}}
		}
		Run(
			WithName("app"),
			WithComponents(
				Provide(dispatcher, di.As(new(Dispatcher))),
			),
		)
		require.True(t, called)
	})

	t.Run("production prefix environment", func(t *testing.T) {
		_ = os.Setenv("ENV", "prod-1")
		called := false
		dispatcher := func(env Env) *testcmp.FuncDispatcher {
			return &testcmp.FuncDispatcher{RunFunc: func(ctx context.Context) error {
				assert.True(t, env.String() == "prod-1")
				assert.False(t, env.IsTest())
				assert.False(t, env.IsDev())
				called = true
				return nil
			}}
		}
		Run(
			WithName("app"),
			WithComponents(
				Provide(dispatcher, di.As(new(Dispatcher))),
			),
		)
		require.True(t, called)
	})

	t.Run("default disabled debug", func(t *testing.T) {
		called := false
		dispatcher := func(info Info) *testcmp.FuncDispatcher {
			return &testcmp.FuncDispatcher{RunFunc: func(ctx context.Context) error {
				assert.False(t, info.Debug)
				called = true
				return nil
			}}
		}
		Run(
			WithName("app"),
			WithComponents(
				Provide(dispatcher, di.As(new(Dispatcher))),
			),
		)
		require.True(t, called)
	})

	t.Run("enable debug by env", func(t *testing.T) {
		_ = os.Setenv("DEBUG", "true")
		called := false
		dispatcher := func(info Info) *testcmp.FuncDispatcher {
			return &testcmp.FuncDispatcher{RunFunc: func(ctx context.Context) error {
				assert.True(t, info.Debug)
				called = true
				return nil
			}}
		}
		Run(
			WithName("app"),
			WithComponents(
				Provide(dispatcher, di.As(new(Dispatcher))),
			),
		)
		require.True(t, called)
	})

}
