package config

import (
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/assert"
)

func TestEnv_Address(t *testing.T) {
	type want struct {
		address              string
		accrualSystemAddress string
		jwtSecret            string
		hasErr               bool
	}

	tests := []struct {
		name string
		env  *env.Options
		want want
	}{
		{
			name: "empty value",
			env: &env.Options{
				Environment: map[string]string{
					"RUN_ADDRESS":            "",
					"ACCRUAL_SYSTEM_ADDRESS": "",
					"JWT_SECRET":             "",
				},
			},
			want: want{
				address:              "",
				accrualSystemAddress: "",
				jwtSecret:            "",
			},
		},
		{
			name: "not empty value",
			env: &env.Options{
				Environment: map[string]string{
					"RUN_ADDRESS":            "RUN_ADDRESS",
					"ACCRUAL_SYSTEM_ADDRESS": "ACCRUAL_SYSTEM_ADDRESS",
					"JWT_SECRET":             "JWT_SECRET",
				},
			},
			want: want{
				address:              "RUN_ADDRESS",
				accrualSystemAddress: "ACCRUAL_SYSTEM_ADDRESS",
				jwtSecret:            "JWT_SECRET",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			envConf, err := NewEnvWithOptions(tt.env)
			if tt.want.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.address, envConf.Address())
				assert.Equal(t, tt.want.jwtSecret, envConf.JWTSecret())
			}
		})
	}
}
