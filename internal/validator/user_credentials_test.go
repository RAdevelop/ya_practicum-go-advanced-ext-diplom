package validator

import (
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/stretchr/testify/assert"
)

func TestIsValidUserCredentials(t *testing.T) {
	type given struct {
		credentials dto.UserCredentials
	}
	tests := []struct {
		name    string
		given   given
		wantErr error
	}{
		{
			name: "valid",
			given: given{
				credentials: dto.UserCredentials{
					Password: "Password",
					Login:    "test",
				},
			},
			wantErr: nil,
		},
		{
			name: "not valid - empty password",
			given: given{
				credentials: dto.UserCredentials{
					Password: "",
					Login:    "test",
				},
			},
			wantErr: perror.ErrInvalidUserCredentials,
		},
		{
			name: "not valid - empty login",
			given: given{
				credentials: dto.UserCredentials{
					Password: "Password",
					Login:    "",
				},
			},
			wantErr: perror.ErrInvalidUserCredentials,
		},
		{
			name: "not valid - empty login",
			given: given{
				credentials: dto.UserCredentials{
					Password: "Password",
					Login:    "   ",
				},
			},
			wantErr: perror.ErrInvalidUserCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantErr, IsValidUserCredentials(tt.given.credentials), "IsValidUserCredentials(%v)", tt.given.credentials)
		})
	}
}
