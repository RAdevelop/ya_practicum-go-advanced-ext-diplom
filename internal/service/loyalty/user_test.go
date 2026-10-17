package loyalty

import (
	"errors"
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func Test_UserLogin(t *testing.T) {

	type given struct {
		inputCustomerCredentials dto.CustomerCredentials
		inputJWTSecret           string
		makeStorage              func(t *testing.T, inputCustomerCredentials dto.CustomerCredentials, inputJWTSecret string, outCustomerModel *model.Customer, outErr error) LoyaltyStorage
	}
	type want struct {
		outCustomerModel *model.Customer
		isTokenEmpty     bool
		outErr           error
	}

	makeStorage := func(t *testing.T, inputCustomerCredentials dto.CustomerCredentials, inputJWTSecret string, outCustomerModel *model.Customer, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().CustomerFindByLogin(t.Context(), inputCustomerCredentials.Login).Return(outCustomerModel, outErr)
		return storage
	}

	login := "TestLogin"
	password := "testPassword"
	passwordHash := hashPassword(t, password)

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success",
			given: given{
				inputCustomerCredentials: dto.CustomerCredentials{
					Login:    login,
					Password: password,
				},
				inputJWTSecret: testJWTSecret,
				makeStorage:    makeStorage,
			},
			want: want{
				outCustomerModel: &model.Customer{
					ID:           1,
					Login:        login,
					PasswordHash: passwordHash,
				},
				isTokenEmpty: false,
				outErr:       nil,
			},
		},
		{
			name: "error ErrCustomerNotFound",
			given: given{
				inputCustomerCredentials: dto.CustomerCredentials{
					Login:    login,
					Password: password,
				},
				inputJWTSecret: testJWTSecret,
				makeStorage:    makeStorage,
			},
			want: want{
				outCustomerModel: nil,
				isTokenEmpty:     true,
				outErr:           perror.ErrCustomerNotFound,
			},
		},
		{
			name: "password Compare error",
			given: given{
				inputCustomerCredentials: dto.CustomerCredentials{
					Login:    login,
					Password: "password Compare error",
				},
				inputJWTSecret: testJWTSecret,
				makeStorage:    makeStorage,
			},
			want: want{
				outCustomerModel: &model.Customer{
					ID:           1,
					Login:        login,
					PasswordHash: passwordHash,
				},
				isTokenEmpty: true,
				outErr:       perror.ErrCustomerInvalidCredentials,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewManager(tt.given.makeStorage(t, tt.given.inputCustomerCredentials, tt.given.inputJWTSecret, tt.want.outCustomerModel, tt.want.outErr))

			token, err := loyaltyManager.UserLogin(t.Context(), tt.given.inputCustomerCredentials, tt.given.inputJWTSecret)
			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)

			if tt.want.isTokenEmpty {
				assert.Emptyf(t, token, "given: %+v", tt.given)
			} else {
				assert.NotEmptyf(t, token, "given: %+v", tt.given)
			}
		})
	}
}

func Test_UserRegister(t *testing.T) {

	type given struct {
		inputCustomerCredentials dto.CustomerCredentials
		inputJWTSecret           string
		makeStorage              func(t *testing.T, inputCustomerCredentials dto.CustomerCredentials, inputJWTSecret string, outCustomerModel *model.Customer, outErr error) LoyaltyStorage
	}
	type want struct {
		outCustomerModel *model.Customer
		isTokenEmpty     bool
		outErr           error
	}

	makeStorage := func(t *testing.T, inputCustomerCredentials dto.CustomerCredentials, inputJWTSecret string, outCustomerModel *model.Customer, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().
			CustomerCreate(
				t.Context(),
				mock.MatchedBy(func(u *model.Customer) bool {
					if u.Login != inputCustomerCredentials.Login {
						return false
					}
					return bcrypt.CompareHashAndPassword(
						[]byte(u.PasswordHash),
						[]byte(inputCustomerCredentials.Password),
					) == nil
				})).
			Return(outCustomerModel, outErr).
			Once()
		return storage
	}

	login := "TestLogin"
	password := "testPassword"
	passwordHash := hashPassword(t, password)

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success",
			given: given{
				inputCustomerCredentials: dto.CustomerCredentials{
					Login:    login,
					Password: password,
				},
				inputJWTSecret: testJWTSecret,
				makeStorage:    makeStorage,
			},
			want: want{
				outCustomerModel: &model.Customer{
					ID:           1,
					Login:        login,
					PasswordHash: passwordHash,
				},
				isTokenEmpty: false,
				outErr:       nil,
			},
		},
		{
			name: "error",
			given: given{
				inputCustomerCredentials: dto.CustomerCredentials{
					Login:    login,
					Password: password,
				},
				inputJWTSecret: testJWTSecret,
				makeStorage:    makeStorage,
			},
			want: want{
				outCustomerModel: nil,
				isTokenEmpty:     true,
				outErr:           errors.New("some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewManager(tt.given.makeStorage(t, tt.given.inputCustomerCredentials, tt.given.inputJWTSecret, tt.want.outCustomerModel, tt.want.outErr))

			token, err := loyaltyManager.UserRegister(t.Context(), tt.given.inputCustomerCredentials, tt.given.inputJWTSecret)
			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)

			if tt.want.isTokenEmpty {
				assert.Emptyf(t, token, "given: %+v", tt.given)
			} else {
				assert.NotEmptyf(t, token, "given: %+v", tt.given)
			}
		})
	}
}
