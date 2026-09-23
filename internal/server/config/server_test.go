package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Address(t *testing.T) {

	type given struct {
		addr string
	}

	type want struct {
		addr string
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "empty address",
			given: given{
				addr: "",
			},
			want: want{
				addr: "",
			},
		},
		{
			name: "ip address",
			given: given{
				addr: "127.0.0.1:8080",
			},
			want: want{
				addr: "127.0.0.1:8080",
			},
		},
		{
			name: "string address",
			given: given{
				addr: "localhost:8080",
			},
			want: want{
				addr: "localhost:8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfgProvider := NewMockProvider(t)
			cfgProvider.EXPECT().Address().Return(tt.want.addr)

			cfg := New(cfgProvider)
			assert.Equal(t, tt.want.addr, cfg.Address())
		})
	}
}
