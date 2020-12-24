package slice

import "testing"

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
