package service

import (
	"errors"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const testJWTSecret = "lPD7WBZ/MCBKK0aEqgzSqfIQSqAGB7VhIfjsZwuXLJE="

var testCustomerDTO = &dto.Customer{
	ID:    1,
	Login: "TestLogin",
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	return string(hash)
}

func Test_OrderUpload(t *testing.T) {

	type given struct {
		inputCustomerDTO *dto.Customer
		inputOrderNumber string

		makeStorage func(t *testing.T, inputCustomerDTO *dto.Customer, inputOrderNumber string, outErr error) LoyaltyStorage
	}
	type want struct {
		outErr error
	}

	makeStorage := func(t *testing.T, inputCustomerDTO *dto.Customer, inputOrderNumber string, outErr error) LoyaltyStorage {

		storage := NewMockLoyaltyStorage(t)
		order := model.Order{
			CustomerID: testCustomerDTO.ID,
			Number:     inputOrderNumber,
			Status:     statusInner.AccrualNew,
		}
		storage.EXPECT().OrderUpload(t.Context(), order).Return(outErr)
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
				inputCustomerDTO: testCustomerDTO,
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
				inputCustomerDTO: testCustomerDTO,
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
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputCustomerDTO, tt.given.inputOrderNumber, tt.want.outErr))

			err := loyaltyManager.OrderUpload(t.Context(), tt.given.inputCustomerDTO, tt.given.inputOrderNumber)

			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)
		})
	}
}

func Test_Orders(t *testing.T) {

	type given struct {
		inputCustomerDTO *dto.Customer
		outOrders        []model.Order
		outErr           error
		makeStorage      func(t *testing.T, inputCustomerDTO *dto.Customer, outOrders []model.Order, outErr error) LoyaltyStorage
	}

	type want struct {
		orders []model.Order
	}

	makeStorage := func(t *testing.T, inputCustomerDTO *dto.Customer, outOrders []model.Order, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().Orders(t.Context(), inputCustomerDTO.ID).Return(outOrders, outErr)
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
				inputCustomerDTO: testCustomerDTO,
				outOrders:        []model.Order{},
				outErr:           nil,
				makeStorage:      makeStorage,
			},
			want: want{
				orders: []model.Order{},
			},
		},
		{
			name: "success and not empty list",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				outOrders: []model.Order{
					{
						Number:     "4532015112830366",
						CustomerID: uint64(1),
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
						CustomerID: 1,
						Status:     statusInner.AccrualNew,
						UploadedAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
					},
				},
			},
		},
		{
			name: "error and nil list",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				outOrders:        nil,
				outErr:           errors.New("test error"),
				makeStorage:      makeStorage,
			},
			want: want{
				orders: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputCustomerDTO, tt.given.outOrders, tt.given.outErr))

			orders, err := loyaltyManager.Orders(t.Context(), tt.given.inputCustomerDTO)

			assert.ErrorIsf(t, err, tt.given.outErr, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.orders, orders, "given: %+v", tt.given)
		})
	}
}

func Test_Balance(t *testing.T) {
	type given struct {
		inputCustomerDTO *dto.Customer
		outBalance       *model.Balance
		outErr           error
		makeStorage      func(t *testing.T, inputCustomerDTO *dto.Customer, outBalance *model.Balance, outErr error) LoyaltyStorage
	}
	type want struct {
		balance *model.Balance
		outErr  error
	}

	makeStorage := func(t *testing.T, inputCustomerDTO *dto.Customer, outBalance *model.Balance, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().Balance(t.Context(), inputCustomerDTO.ID).Return(outBalance, outErr)
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
				inputCustomerDTO: testCustomerDTO,
				outBalance:       &model.Balance{},
				makeStorage:      makeStorage,
			},
			want: want{
				balance: &model.Balance{},
				outErr:  nil,
			},
		},
		{
			name: "success and nil balance",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				outBalance:       nil,
				makeStorage:      makeStorage,
			},
			want: want{
				balance: nil,
				outErr:  nil,
			},
		},
		{
			name: "success and not empty balance",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				outBalance: &model.Balance{
					ID:         1,
					CustomerID: 1,
					Current:    19.05,
					Withdrawn:  2019.05,
				},
				makeStorage: makeStorage,
			},
			want: want{
				balance: &model.Balance{
					ID:         1,
					CustomerID: 1,
					Current:    19.05,
					Withdrawn:  2019.05,
				},
				outErr: nil,
			},
		},
		{
			name: "error and empty balance",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				outBalance:       nil,
				makeStorage:      makeStorage,
			},
			want: want{
				balance: nil,
				outErr:  errors.New("test error"),
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputCustomerDTO, tt.given.outBalance, tt.want.outErr))

			balance, err := loyaltyManager.Balance(t.Context(), tt.given.inputCustomerDTO)

			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.balance, balance, "given: %+v", tt.given)
		})
	}
}

func Test_BalanceWithdrawals(t *testing.T) {
	type given struct {
		inputCustomerDTO *dto.Customer
		outWithdrawals   []model.Withdrawal
		outErr           error
		makeStorage      func(t *testing.T, inputCustomerDTO *dto.Customer, outWithdrawals []model.Withdrawal, outErr error) LoyaltyStorage
	}
	type want struct {
		withdrawals []model.Withdrawal
		outErr      error
	}

	makeStorage := func(t *testing.T, inputCustomerDTO *dto.Customer, outWithdrawals []model.Withdrawal, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().BalanceWithdrawals(t.Context(), inputCustomerDTO.ID).Return(outWithdrawals, outErr)
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
				inputCustomerDTO: testCustomerDTO,
				outWithdrawals:   nil,
				outErr:           nil,
				makeStorage:      makeStorage,
			},
			want: want{
				withdrawals: nil,
			},
		},
		{
			name: "error and empty list",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				outWithdrawals:   nil,

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
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputCustomerDTO, tt.given.outWithdrawals, tt.want.outErr))

			withdrawals, err := loyaltyManager.BalanceWithdrawals(t.Context(), tt.given.inputCustomerDTO)

			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.withdrawals, withdrawals, "given: %+v", tt.given)
		})
	}
}

func Test_BalanceWithdraw(t *testing.T) {
	type given struct {
		inputCustomerDTO *dto.Customer
		inputOrderNumber string
		inputSum         float64
		makeStorage      func(t *testing.T, inputCustomerDTO *dto.Customer, inputOrderNumber string, inputSum float64, outErr error) LoyaltyStorage
	}

	type want struct {
		outErr error
	}

	makeStorage := func(t *testing.T, inputCustomerDTO *dto.Customer, inputOrderNumber string, inputSum float64, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().BalanceWithdraw(t.Context(), inputCustomerDTO.ID, inputOrderNumber, inputSum).Return(outErr)

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
				inputCustomerDTO: testCustomerDTO,
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
				inputCustomerDTO: testCustomerDTO,
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
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputCustomerDTO, tt.given.inputOrderNumber, tt.given.inputSum, tt.want.outErr))

			err := loyaltyManager.BalanceWithdraw(t.Context(), tt.given.inputCustomerDTO, tt.given.inputOrderNumber, tt.given.inputSum)

			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)
		})
	}
}

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
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputCustomerCredentials, tt.given.inputJWTSecret, tt.want.outCustomerModel, tt.want.outErr))

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
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputCustomerCredentials, tt.given.inputJWTSecret, tt.want.outCustomerModel, tt.want.outErr))

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
