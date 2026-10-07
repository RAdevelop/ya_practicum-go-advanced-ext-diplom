package service

import (
	"errors"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_Balance(t *testing.T) {
	type given struct {
		inputCustomerDTO *dto.Customer
		outBalance       *model.Balance
		outErr           error
		makeStorage      func(t *testing.T, inputCustomerDTO *dto.Customer, outBalance *model.Balance, outErr error) LoyaltyStorage
	}
	type want struct {
		balance dto.Balance
		outErr  error
	}

	makeStorage := func(t *testing.T, inputCustomerDTO *dto.Customer, outBalance *model.Balance, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().BalanceByCustomerID(t.Context(), inputCustomerDTO.ID).Return(outBalance, outErr)
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
				balance: dto.Balance{},
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
				balance: dto.Balance{
					Current:   19.05,
					Withdrawn: 2019.05,
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
				balance: dto.Balance{},
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
		withdrawals []dto.Withdrawal
		err         error
	}

	makeStorage := func(t *testing.T, inputCustomerDTO *dto.Customer, outWithdrawals []model.Withdrawal, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().BalanceWithdrawalsByCustomerID(t.Context(), inputCustomerDTO.ID).Return(outWithdrawals, outErr)
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
				outWithdrawals: []model.Withdrawal{
					{
						ID:          1,
						CustomerID:  1,
						Order:       "4532015112830366",
						Sum:         19.0,
						ProcessedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				outErr:      nil,
				makeStorage: makeStorage,
			},
			want: want{
				withdrawals: []dto.Withdrawal{
					{
						Order:       "4532015112830366",
						Sum:         19.0,
						ProcessedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				err: nil,
			},
		},
		{
			name: "error and empty list",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				outWithdrawals:   nil,
				outErr:           errors.New("test error"),
				makeStorage:      makeStorage,
			},
			want: want{
				withdrawals: nil,
				err:         errors.New("test error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewLoyaltyManager(tt.given.makeStorage(t, tt.given.inputCustomerDTO, tt.given.outWithdrawals, tt.want.err))

			withdrawals, err := loyaltyManager.BalanceWithdrawals(t.Context(), tt.given.inputCustomerDTO)

			assert.ErrorIsf(t, err, tt.want.err, "given: %+v", tt.given)
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
		withdrawal := model.Withdrawal{
			CustomerID: inputCustomerDTO.ID,
			Order:      inputOrderNumber,
			Sum:        inputSum,
		}
		storage.EXPECT().BalanceWithdraw(t.Context(), withdrawal).Return(outErr)

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
