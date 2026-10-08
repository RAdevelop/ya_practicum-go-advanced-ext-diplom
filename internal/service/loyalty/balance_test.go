package loyalty

import (
	"errors"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	statusOuter "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/outer"
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
			loyaltyManager := NewManager(tt.given.makeStorage(t, tt.given.inputCustomerDTO, tt.given.outBalance, tt.want.outErr))

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
						OrderID:     1,
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
			loyaltyManager := NewManager(tt.given.makeStorage(t, tt.given.inputCustomerDTO, tt.given.outWithdrawals, tt.want.err))

			withdrawals, err := loyaltyManager.BalanceWithdrawals(t.Context(), tt.given.inputCustomerDTO)

			assert.ErrorIsf(t, err, tt.want.err, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.withdrawals, withdrawals, "given: %+v", tt.given)
		})
	}
}

func Test_BalanceWithdraw(t *testing.T) {
	type given struct {
		inputCustomerDTO *dto.Customer
		inputWithdraw    dto.BalanceWithdraw
		makeStorage      func(t *testing.T, inputCustomerDTO *dto.Customer, inputWithdraw dto.BalanceWithdraw, outErr error) LoyaltyStorage
	}

	type want struct {
		outErr error
	}

	makeStorage := func(t *testing.T, inputCustomerDTO *dto.Customer, inputWithdraw dto.BalanceWithdraw, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		withdrawal := model.Withdrawal{
			CustomerID: inputCustomerDTO.ID,
			Order:      inputWithdraw.OrderNumber,
			Sum:        inputWithdraw.Sum,
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
				inputWithdraw: dto.BalanceWithdraw{
					OrderNumber: "4532015112830366",
					Sum:         19.0,
				},
				makeStorage: makeStorage,
			},
			want: want{
				outErr: nil,
			},
		},
		{
			name: "error",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				inputWithdraw: dto.BalanceWithdraw{
					OrderNumber: "4532015112830366",
					Sum:         19.0,
				},
				makeStorage: makeStorage,
			},
			want: want{
				outErr: errors.New("test error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewManager(tt.given.makeStorage(t, tt.given.inputCustomerDTO, tt.given.inputWithdraw, tt.want.outErr))

			err := loyaltyManager.BalanceWithdraw(t.Context(), tt.given.inputCustomerDTO, tt.given.inputWithdraw)

			assert.ErrorIsf(t, err, tt.want.outErr, "given: %+v", tt.given)
		})
	}
}

func Test_BalanceAccrual(t *testing.T) {

	type given struct {
		inputAccrual dto.Accrual
		outErr       error
		makeStorage  func(t *testing.T, inputAccrual dto.Accrual, outErr error) LoyaltyStorage
	}
	type want struct {
		err error
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success",
			given: given{
				inputAccrual: dto.Accrual{
					Order:   "4532015112830366",
					Status:  statusOuter.AccrualProcessed,
					Accrual: 111.0,
				},
				outErr: nil,
				makeStorage: func(t *testing.T, inputAccrual dto.Accrual, outErr error) LoyaltyStorage {
					storage := NewMockLoyaltyStorage(t)

					accrualModel := model.Accrual{
						Order:   inputAccrual.Order,
						Status:  inputAccrual.Status,
						Accrual: inputAccrual.Accrual,
					}
					storage.EXPECT().BalanceAccrual(t.Context(), accrualModel).Return(outErr)
					return storage
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "ErrAccrualApply",
			given: given{
				inputAccrual: dto.Accrual{
					Order:  "4532015112830366",
					Status: statusOuter.AccrualProcessing,
				},
				outErr: perror.ErrAccrualApply,
				makeStorage: func(t *testing.T, inputAccrual dto.Accrual, outErr error) LoyaltyStorage {
					storage := NewMockLoyaltyStorage(t)

					accrualModel := model.Accrual{
						Order:  inputAccrual.Order,
						Status: inputAccrual.Status,
					}
					storage.EXPECT().BalanceAccrual(t.Context(), accrualModel).Return(outErr)
					return storage
				},
			},
			want: want{
				err: perror.ErrAccrualApply,
			},
		},
		{
			name: "ErrAccrualAlreadyProcessed",
			given: given{
				inputAccrual: dto.Accrual{
					Order:  "4532015112830366",
					Status: statusOuter.AccrualProcessing,
				},
				outErr: perror.ErrOrderAccrualAlreadyProcessed,
				makeStorage: func(t *testing.T, inputAccrual dto.Accrual, outErr error) LoyaltyStorage {
					storage := NewMockLoyaltyStorage(t)

					accrualModel := model.Accrual{
						Order:  inputAccrual.Order,
						Status: inputAccrual.Status,
					}
					storage.EXPECT().BalanceAccrual(t.Context(), accrualModel).Return(outErr)
					return storage
				},
			},
			want: want{
				err: perror.ErrOrderAccrualAlreadyProcessed,
			},
		},
		{
			name: "ErrBalanceIncrement",
			given: given{
				inputAccrual: dto.Accrual{
					Order:  "4532015112830366",
					Status: statusOuter.AccrualProcessing,
				},
				outErr: perror.ErrBalanceIncrement,
				makeStorage: func(t *testing.T, inputAccrual dto.Accrual, outErr error) LoyaltyStorage {
					storage := NewMockLoyaltyStorage(t)

					accrualModel := model.Accrual{
						Order:  inputAccrual.Order,
						Status: inputAccrual.Status,
					}
					storage.EXPECT().BalanceAccrual(t.Context(), accrualModel).Return(outErr)
					return storage
				},
			},
			want: want{
				err: perror.ErrBalanceIncrement,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewManager(tt.given.makeStorage(t, tt.given.inputAccrual, tt.given.outErr))
			err := loyaltyManager.BalanceAccrual(t.Context(), tt.given.inputAccrual)

			assert.ErrorIsf(t, err, tt.want.err, "given: %+v", tt.given)
		})
	}
}
