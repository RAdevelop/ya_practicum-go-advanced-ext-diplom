package service

import (
	"errors"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/jwtoken"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const testJWTSecret = "lPD7WBZ/MCBKK0aEqgzSqfIQSqAGB7VhIfjsZwuXLJE="

var testUserDTO = &dto.User{
	ID:    1,
	Login: "TestLogin",
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	return string(hash)
}

// tokenGenerate - генерируем токен, если userDTO задан
func tokenGenerate(t *testing.T, userDTO *dto.User) string {
	t.Helper()

	if userDTO != nil {
		token, err := jwtoken.Generate(*userDTO, []byte(testJWTSecret))
		assert.NoErrorf(t, err, "userDTO: %+v", userDTO)
		return token
	}

	return ""
}

func Test_OrderUpload(t *testing.T) {

	type given struct {
		inputUserDTO     *dto.User
		inputOrderNumber string

		makeStorage func(t *testing.T, inputUserDTO *dto.User, inputOrderNumber string, outErr error) LoyaltyStorage
	}
	type want struct {
		outErr error
	}

	makeStorage := func(t *testing.T, inputUserDTO *dto.User, inputOrderNumber string, outErr error) LoyaltyStorage {

		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().OrderUpload(t.Context(), inputUserDTO.ID, inputOrderNumber).Return(outErr)
		return storage
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success",
			given: given{
				inputUserDTO:     testUserDTO,
				inputOrderNumber: "1",
				makeStorage:      makeStorage,
			},
			want: want{
				outErr: nil,
			},
		},
		{
			name: "error",
			given: given{
				inputUserDTO:     testUserDTO,
				inputOrderNumber: "1",
				makeStorage:      makeStorage,
			},
			want: want{
				outErr: errors.New("test error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputUserDTO, tt.given.inputOrderNumber, tt.want.outErr))

			err := loyaltyManager.OrderUpload(t.Context(), tt.given.inputUserDTO, tt.given.inputOrderNumber)

			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)
		})
	}
}

func Test_Orders(t *testing.T) {

	type given struct {
		inputUserDTO *dto.User
		outOrders    []model.Order
		outErr       error
		makeStorage  func(t *testing.T, inputUserDTO *dto.User, outOrders []model.Order, outErr error) LoyaltyStorage
	}

	type want struct {
		orders []model.Order
	}

	makeStorage := func(t *testing.T, inputUserDTO *dto.User, outOrders []model.Order, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().Orders(t.Context(), inputUserDTO.ID).Return(outOrders, outErr)
		return storage
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success and empty list",
			given: given{
				inputUserDTO: testUserDTO,
				outOrders:    []model.Order{},
				outErr:       nil,
				makeStorage:  makeStorage,
			},
			want: want{
				orders: []model.Order{},
			},
		},
		{
			name: "success and not empty list",
			given: given{
				inputUserDTO: testUserDTO,
				outOrders: []model.Order{
					{
						Number:     "4532015112830366",
						UserID:     1,
						Status:     statusInner.AccrualNew,
						UploadedAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
					},
				},
				outErr:      nil,
				makeStorage: makeStorage,
			},
			want: want{
				orders: []model.Order{
					{
						Number:     "4532015112830366",
						UserID:     1,
						Status:     statusInner.AccrualNew,
						UploadedAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
					},
				},
			},
		},
		{
			name: "error and nil list",
			given: given{
				inputUserDTO: testUserDTO,
				outOrders:    nil,
				outErr:       errors.New("test error"),
				makeStorage:  makeStorage,
			},
			want: want{
				orders: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputUserDTO, tt.given.outOrders, tt.given.outErr))

			orders, err := loyaltyManager.Orders(t.Context(), tt.given.inputUserDTO)

			assert.ErrorIsf(t, err, tt.given.outErr, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.orders, orders, "given: %+v", tt.given)
		})
	}
}

func Test_Balance(t *testing.T) {
	type given struct {
		inputUserDTO *dto.User
		outBalance   *model.Balance
		outErr       error
		makeStorage  func(t *testing.T, inputUserDTO *dto.User, outBalance *model.Balance, outErr error) LoyaltyStorage
	}
	type want struct {
		balance *model.Balance
		outErr  error
	}

	makeStorage := func(t *testing.T, inputUserDTO *dto.User, outBalance *model.Balance, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().Balance(t.Context(), inputUserDTO.ID).Return(outBalance, outErr)
		return storage
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success and empty balance",
			given: given{
				inputUserDTO: testUserDTO,
				outBalance:   &model.Balance{},
				makeStorage:  makeStorage,
			},
			want: want{
				balance: &model.Balance{},
				outErr:  nil,
			},
		},
		{
			name: "success and nil balance",
			given: given{
				inputUserDTO: testUserDTO,
				outBalance:   nil,
				makeStorage:  makeStorage,
			},
			want: want{
				balance: nil,
				outErr:  nil,
			},
		},
		{
			name: "success and not empty balance",
			given: given{
				inputUserDTO: testUserDTO,
				outBalance: &model.Balance{
					ID:        1,
					UserID:    1,
					Current:   19.05,
					Withdrawn: 2019.05,
				},
				makeStorage: makeStorage,
			},
			want: want{
				balance: &model.Balance{
					ID:        1,
					UserID:    1,
					Current:   19.05,
					Withdrawn: 2019.05,
				},
				outErr: nil,
			},
		},
		{
			name: "error and empty balance",
			given: given{
				inputUserDTO: testUserDTO,
				outBalance:   nil,
				makeStorage:  makeStorage,
			},
			want: want{
				balance: nil,
				outErr:  errors.New("test error"),
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputUserDTO, tt.given.outBalance, tt.want.outErr))

			balance, err := loyaltyManager.Balance(t.Context(), tt.given.inputUserDTO)

			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.balance, balance, "given: %+v", tt.given)
		})
	}
}

func Test_BalanceWithdrawals(t *testing.T) {
	type given struct {
		inputUserDTO   *dto.User
		outWithdrawals []model.Withdrawal
		outErr         error
		makeStorage    func(t *testing.T, inputUserDTO *dto.User, outWithdrawals []model.Withdrawal, outErr error) LoyaltyStorage
	}
	type want struct {
		withdrawals []model.Withdrawal
		outErr      error
	}

	makeStorage := func(t *testing.T, inputUserDTO *dto.User, outWithdrawals []model.Withdrawal, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().BalanceWithdrawals(t.Context(), inputUserDTO.ID).Return(outWithdrawals, outErr)
		return storage
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success and empty list",
			given: given{
				inputUserDTO:   testUserDTO,
				outWithdrawals: nil,
				outErr:         nil,
				makeStorage:    makeStorage,
			},
			want: want{
				withdrawals: nil,
			},
		},
		{
			name: "error and empty list",
			given: given{
				inputUserDTO:   testUserDTO,
				outWithdrawals: nil,

				makeStorage: makeStorage,
			},
			want: want{
				withdrawals: nil,
				outErr:      errors.New("test error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputUserDTO, tt.given.outWithdrawals, tt.want.outErr))

			withdrawals, err := loyaltyManager.BalanceWithdrawals(t.Context(), tt.given.inputUserDTO)

			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.withdrawals, withdrawals, "given: %+v", tt.given)
		})
	}
}

func Test_BalanceWithdraw(t *testing.T) {
	type given struct {
		inputUserDTO     *dto.User
		inputOrderNumber string
		inputSum         float64
		makeStorage      func(t *testing.T, inputUserDTO *dto.User, inputOrderNumber string, inputSum float64, outErr error) LoyaltyStorage
	}

	type want struct {
		outErr error
	}

	makeStorage := func(t *testing.T, inputUserDTO *dto.User, inputOrderNumber string, inputSum float64, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().BalanceWithdraw(t.Context(), inputUserDTO.ID, inputOrderNumber, inputSum).Return(outErr)

		return storage
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success",
			given: given{
				inputUserDTO:     testUserDTO,
				inputOrderNumber: "4532015112830366",
				inputSum:         19.0,
				makeStorage:      makeStorage,
			},
			want: want{
				outErr: nil,
			},
		},
		{
			name: "error",
			given: given{
				inputUserDTO:     testUserDTO,
				inputOrderNumber: "4532015112830366",
				inputSum:         19.0,
				makeStorage:      makeStorage,
			},
			want: want{
				outErr: errors.New("test error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputUserDTO, tt.given.inputOrderNumber, tt.given.inputSum, tt.want.outErr))

			err := loyaltyManager.BalanceWithdraw(t.Context(), tt.given.inputUserDTO, tt.given.inputOrderNumber, tt.given.inputSum)

			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)
		})
	}
}

func Test_UserLogin(t *testing.T) {

	type given struct {
		inputUserCredentials dto.UserCredentials
		inputJWTSecret       string
		makeStorage          func(t *testing.T, inputUserCredentials dto.UserCredentials, inputJWTSecret string, outUserModel *model.User, outErr error) LoyaltyStorage
	}
	type want struct {
		outUserModel *model.User
		isTokenEmpty bool
		outErr       error
	}

	makeStorage := func(t *testing.T, inputUserCredentials dto.UserCredentials, inputJWTSecret string, outUserModel *model.User, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().UserByLogin(t.Context(), inputUserCredentials.Login).Return(outUserModel, outErr)
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
				inputUserCredentials: dto.UserCredentials{
					Login:    login,
					Password: password,
				},
				inputJWTSecret: testJWTSecret,
				makeStorage:    makeStorage,
			},
			want: want{
				outUserModel: &model.User{
					ID:           1,
					Login:        login,
					PasswordHash: passwordHash,
				},
				isTokenEmpty: false,
				outErr:       nil,
			},
		},
		{
			name: "error ErrUserNotFound",
			given: given{
				inputUserCredentials: dto.UserCredentials{
					Login:    login,
					Password: password,
				},
				inputJWTSecret: testJWTSecret,
				makeStorage:    makeStorage,
			},
			want: want{
				outUserModel: nil,
				isTokenEmpty: true,
				outErr:       perror.ErrUserNotFound,
			},
		},
		{
			name: "password Compare error",
			given: given{
				inputUserCredentials: dto.UserCredentials{
					Login:    login,
					Password: "password Compare error",
				},
				inputJWTSecret: testJWTSecret,
				makeStorage:    makeStorage,
			},
			want: want{
				outUserModel: &model.User{
					ID:           1,
					Login:        login,
					PasswordHash: passwordHash,
				},
				isTokenEmpty: true,
				outErr:       perror.ErrInvalidUserCredentials,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputUserCredentials, tt.given.inputJWTSecret, tt.want.outUserModel, tt.want.outErr))

			token, err := loyaltyManager.UserLogin(t.Context(), tt.given.inputUserCredentials, tt.given.inputJWTSecret)
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
		inputUserCredentials dto.UserCredentials
		inputJWTSecret       string
		makeStorage          func(t *testing.T, inputUserCredentials dto.UserCredentials, inputJWTSecret string, outUserModel *model.User, outErr error) LoyaltyStorage
	}
	type want struct {
		outUserModel *model.User
		isTokenEmpty bool
		outErr       error
	}

	makeStorage := func(t *testing.T, inputUserCredentials dto.UserCredentials, inputJWTSecret string, outUserModel *model.User, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().
			UserCreate(
				t.Context(),
				mock.MatchedBy(func(u *model.User) bool {
					if u.Login != inputUserCredentials.Login {
						return false
					}
					return bcrypt.CompareHashAndPassword(
						[]byte(u.PasswordHash),
						[]byte(inputUserCredentials.Password),
					) == nil
				})).
			Return(outUserModel, outErr).
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
				inputUserCredentials: dto.UserCredentials{
					Login:    login,
					Password: password,
				},
				inputJWTSecret: testJWTSecret,
				makeStorage:    makeStorage,
			},
			want: want{
				outUserModel: &model.User{
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
				inputUserCredentials: dto.UserCredentials{
					Login:    login,
					Password: password,
				},
				inputJWTSecret: testJWTSecret,
				makeStorage:    makeStorage,
			},
			want: want{
				outUserModel: nil,
				isTokenEmpty: true,
				outErr:       errors.New("some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputUserCredentials, tt.given.inputJWTSecret, tt.want.outUserModel, tt.want.outErr))

			token, err := loyaltyManager.UserRegister(t.Context(), tt.given.inputUserCredentials, tt.given.inputJWTSecret)
			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)

			if tt.want.isTokenEmpty {
				assert.Emptyf(t, token, "given: %+v", tt.given)
			} else {
				assert.NotEmptyf(t, token, "given: %+v", tt.given)
			}
		})
	}
}
