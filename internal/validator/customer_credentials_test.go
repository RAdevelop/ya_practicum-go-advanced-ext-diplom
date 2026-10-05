package validator

import (
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/stretchr/testify/assert"
)

func TestIsValidCustomerCredentials(t *testing.T) {
	type given struct {
		credentials dto.CustomerCredentials
	}
	tests := []struct {
		name    string
		given   given
		wantErr error
	}{
		{
			name: "valid",
			given: given{
				credentials: dto.CustomerCredentials{
					Password: "Password",
					Login:    "test",
				},
			},
			wantErr: nil,
		},
		{
			name: "not valid - empty password",
			given: given{
				credentials: dto.CustomerCredentials{
					Password: "",
					Login:    "test",
				},
			},
			wantErr: perror.ErrInvalidCustomerCredentials,
		},
		{
			name: "not valid - empty login",
			given: given{
				credentials: dto.CustomerCredentials{
					Password: "Password",
					Login:    "",
				},
			},
			wantErr: perror.ErrInvalidCustomerCredentials,
		},
		{
			name: "not valid - empty login",
			given: given{
				credentials: dto.CustomerCredentials{
					Password: "Password",
					Login:    "   ",
				},
			},
			wantErr: perror.ErrInvalidCustomerCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantErr, IsValidCustomerCredentials(tt.given.credentials), "IsValidCustomerCredentials(%v)", tt.given.credentials)
		})
	}
}
