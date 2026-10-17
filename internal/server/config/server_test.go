package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Address(t *testing.T) {

	type given struct {
		addr                 string
		accrualSystemAddress string
		jwtSecret            string
	}

	type want struct {
		addr                 string
		accrualSystemAddress string
		jwtSecret            string
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "empty value",
			given: given{
				addr:                 "",
				accrualSystemAddress: "",
				jwtSecret:            "",
			},
			want: want{
				addr:                 "",
				accrualSystemAddress: "",
				jwtSecret:            "",
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
				addr:                 "localhost:8080",
				accrualSystemAddress: "localhost:8081",
				jwtSecret:            "jwtSecret",
			},
			want: want{
				addr:                 "localhost:8080",
				accrualSystemAddress: "localhost:8081",
				jwtSecret:            "jwtSecret",
			},
		},
		{
			name: "jwtSecret",
			given: given{
				jwtSecret: "jwtSecret",
			},
			want: want{
				jwtSecret: "jwtSecret",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfgProvider := NewMockProvider(t)
			cfgProvider.EXPECT().Address().Return(tt.want.addr)
			cfgProvider.EXPECT().AccrualSystemAddress().Return(tt.want.accrualSystemAddress)
			cfgProvider.EXPECT().JWTSecret().Return(tt.want.jwtSecret)

			cfg := New(cfgProvider)
			assert.Equal(t, tt.want.addr, cfg.Address())
			assert.Equal(t, tt.want.accrualSystemAddress, cfg.AccrualSystemAddress())
			assert.Equal(t, tt.want.jwtSecret, cfg.JWTSecret())
		})
	}
}
