package slice

import "testing"

func TestInfo_IsDev(t *testing.T) {
	type fields struct {
		Name  string
		Env   Env
		Debug bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Dev mode",
			fields: fields{
				Name:  "",
				Env:   "dev",
				Debug: false,
			},
			want: true,
		},
		{
			name: "not Dev mode",
			fields: fields{
				Name:  "",
				Env:   "",
				Debug: false,
			},
			want: false,
		},
		{
			name: "Dev mode with uppercase env",
			fields: fields{
				Name:  "",
				Env:   "DEV",
				Debug: false,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &Info{
				Name:  tt.fields.Name,
				Env:   tt.fields.Env,
				Debug: tt.fields.Debug,
			}
			if got := info.IsDev(); got != tt.want {
				t.Errorf("IsDev() = %v, want %v", got, tt.want)
			}
		})
	}
}
