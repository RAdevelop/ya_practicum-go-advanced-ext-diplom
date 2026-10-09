package database

import (
	"context"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestStorage_BalanceByCustomerID(t *testing.T) {

	type given struct {
		customerID   uint64
		ordersCreate []model.Order
		accruals     []model.Accrual
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
			name: "ErrCustomerInvalidCredentials",
			given: given{
				customerID:   0,
				ordersCreate: []model.Order{},
				accruals:     []model.Accrual{},
			},
			want: want{
				err: perror.ErrCustomerInvalidCredentials,
			},
		},
		{
			name: "success",
			given: given{
				customerID: 1,
				ordersCreate: []model.Order{
					{
						Number:     "4532015112830366",
						CustomerID: uint64(1),
						Status:     statusInner.AccrualNew,
					},
					{
						Number:     "4532015112830368",
						CustomerID: uint64(1),
						Status:     statusInner.AccrualProcessing,
					},
				},
				accruals: []model.Accrual{
					{
						Order:   "4532015112830366",
						Status:  statusInner.AccrualProcessed,
						Accrual: 20,
					},
					{
						Order:   "4532015112830368",
						Status:  statusInner.AccrualProcessed,
						Accrual: 30,
					},
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	storage := setUpStorage(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.DB.RunInTransaction(t.Context(), func(ctx context.Context) error {

				//создать заказ
				for _, order := range tt.given.ordersCreate {
					err := storage.OrderUpload(ctx, order)
					assert.NoErrorf(t, err, "given:%+v", tt.given)
				}

				//начислить баллы
				for _, accrual := range tt.given.accruals {
					err := storage.BalanceAccrual(ctx, accrual)
					assert.NoErrorf(t, err, "given:%+v", tt.given)
				}

				//получить баланс
				balance, err := storage.BalanceByCustomerID(ctx, tt.given.customerID)
				assert.ErrorIsf(t, err, tt.want.err, "given:%+v", tt.given)

				if tt.want.err == nil {
					var accrualsSum float64
					for _, accrual := range tt.given.accruals {
						accrualsSum += accrual.Accrual
					}

					checkBalance(t, tt.given.customerID, accrualsSum, 0, balance, tt.given)
				}

				return errForTransactionRollback

			}, pgx.TxOptions{IsoLevel: pgx.Serializable})

			assert.Errorf(t, err, "given:%+v", tt.given)
		})
	}
}

func TestStorage_BalanceAccrual(t *testing.T) {

	type given struct {
		orderCreate model.Order
		accrual     model.Accrual
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
			name: "ErrOrderNotFound",
			given: given{
				orderCreate: model.Order{
					Number:     "4532015112830366",
					CustomerID: uint64(1),
					Status:     statusInner.AccrualNew,
				},
				accrual: model.Accrual{
					Order:   "4532015112830369", //нет такого заказа
					Status:  statusInner.AccrualProcessed,
					Accrual: 20,
				},
			},
			want: want{
				err: perror.ErrOrderNotFound,
			},
		},
		{
			name: "ErrOrderAccrualAlreadyProcessed",
			given: given{
				orderCreate: model.Order{
					Number:     "4532015112830366",
					CustomerID: uint64(1),
					Status:     statusInner.AccrualProcessed,
				},
				accrual: model.Accrual{
					Order:   "4532015112830366",
					Status:  statusInner.AccrualProcessed,
					Accrual: 20,
				},
			},
			want: want{
				err: perror.ErrOrderAccrualAlreadyProcessed,
			},
		},
		{
			name: "success",
			given: given{
				orderCreate: model.Order{
					Number:     "4532015112830366",
					CustomerID: uint64(1),
					Status:     statusInner.AccrualProcessing,
				},
				accrual: model.Accrual{
					Order:   "4532015112830366",
					Status:  statusInner.AccrualProcessed,
					Accrual: 20,
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	storage := setUpStorage(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := storage.DB.RunInTransaction(t.Context(), func(ctx context.Context) error {

				err := storage.OrderUpload(ctx, tt.given.orderCreate)
				assert.NoError(t, err, "given:%+v", tt.given)

				err = storage.BalanceAccrual(ctx, tt.given.accrual)
				assert.ErrorIs(t, err, tt.want.err, "given:%+v", tt.given)

				return errForTransactionRollback
			}, pgx.TxOptions{IsoLevel: pgx.Serializable})

			assert.Error(t, err, "given:%+v", tt.given)
		})
	}
}

func TestStorage_BalanceWithdrawalsByCustomerID(t *testing.T) {
	//TODO implement
	t.Skip("TODO implement")
}
func TestStorage_BalanceWithdraw(t *testing.T) {

	type want struct {
		withdrawal  model.Withdrawal
		expectedErr error
	}

	type given struct {
		customer     *model.Customer
		ordersCreate []model.Order
		accruals     []model.Accrual
		withdrawals  []want
	}

	customerModel := &model.Customer{
		Login:        "john",
		PasswordHash: "PasswordHash",
	}

	tests := []struct {
		name  string
		given given
	}{
		{
			name: "ErrWithdrawalInvalid",
			given: given{
				customer:     customerModel,
				ordersCreate: []model.Order{},
				accruals:     []model.Accrual{},
				withdrawals: []want{
					{
						withdrawal:  model.Withdrawal{},
						expectedErr: perror.ErrWithdrawalInvalid,
					},
				},
			},
		},
		{
			name: "ErrOrderNotFound",
			given: given{
				customer: customerModel,
				ordersCreate: []model.Order{
					{
						Number: "4532015112830366",
						//CustomerID: для первого заказа в таком списке в тесте ниже будет подставлен id созданного покупателя,
						Status: statusInner.AccrualNew,
					},
				},
				accruals: []model.Accrual{
					{
						Order:   "4532015112830366",
						Status:  statusInner.AccrualProcessed,
						Accrual: 20,
					},
				},
				withdrawals: []want{
					{
						withdrawal: model.Withdrawal{
							//CustomerID: для всех списаний в таком списке в тесте ниже будет подставлен id созданного покупателя,
							Order: "4532015112830368", // ErrOrderNotFound ТАКОГО ЗАКАЗА НЕТ У ПОКУПАТЕЛЯ
							Sum:   20,
						},
						expectedErr: perror.ErrOrderNotFound,
					},
				},
			},
		},
		{
			name: "ErrBalanceInsufficient",
			given: given{
				customer: customerModel,
				ordersCreate: []model.Order{
					{
						Number: "4532015112830366",
						//CustomerID: для первого заказа в таком списке в тесте ниже будет подставлен id созданного покупателя,
						Status: statusInner.AccrualNew,
					},
				},
				accruals: []model.Accrual{
					{
						Order:   "4532015112830366",
						Status:  statusInner.AccrualProcessed,
						Accrual: 20,
					},
				},
				withdrawals: []want{
					{
						withdrawal: model.Withdrawal{
							//CustomerID: для всех списаний в таком списке в тесте ниже будет подставлен id созданного покупателя,
							Order: "4532015112830366",
							Sum:   21, // сумма списания больше чем есть на балансе
						},
						expectedErr: perror.ErrBalanceInsufficient,
					},
				},
			},
		},
		{
			name: "ErrBalanceInsufficient",
			given: given{
				customer: customerModel,
				ordersCreate: []model.Order{
					{
						Number: "4532015112830366",
						//CustomerID: для первого заказа в таком списке в тесте ниже будет подставлен id созданного покупателя,
						Status: statusInner.AccrualNew,
					},
				},
				accruals: []model.Accrual{
					{
						Order:   "4532015112830366",
						Status:  statusInner.AccrualProcessed,
						Accrual: 20,
					},
				},
				withdrawals: []want{
					{
						withdrawal: model.Withdrawal{
							//CustomerID: для всех списаний в таком списке в тесте ниже будет подставлен id созданного покупателя,
							Order: "4532015112830366",
							Sum:   20,
						},
						expectedErr: nil,
					},
					{
						withdrawal: model.Withdrawal{
							//CustomerID: для всех списаний в таком списке в тесте ниже будет подставлен id созданного покупателя,
							Order: "4532015112830366",
							Sum:   20,
						},
						expectedErr: perror.ErrBalanceInsufficient, //ВТОРОЕ СПИСАНИЕ ПО ТОМУ ЖЕ ЗАКАЗУ
					},
				},
			},
		},
		{
			name: "success",
			given: given{
				customer: customerModel,
				ordersCreate: []model.Order{
					{
						Number: "4532015112830366",
						//CustomerID: для первого заказа в таком списке в тесте ниже будет подставлен id созданного покупателя,
						Status: statusInner.AccrualNew,
					},
				},
				accruals: []model.Accrual{
					{
						Order:   "4532015112830366",
						Status:  statusInner.AccrualProcessed,
						Accrual: 20,
					},
				},
				withdrawals: []want{
					{
						withdrawal: model.Withdrawal{
							//CustomerID: для всех списаний в таком списке в тесте ниже будет подставлен id созданного покупателя,
							Order: "4532015112830366",
							Sum:   20,
						},
						expectedErr: nil,
					},
				},
			},
		},
	}

	storage := setUpStorage(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.DB.RunInTransaction(t.Context(), func(ctx context.Context) error {

				//создать покупателя
				customer, err := storage.CustomerCreate(ctx, tt.given.customer)
				assert.NoErrorf(t, err, "given:%+v", tt.given)

				//создать заказ(ы)
				for i, order := range tt.given.ordersCreate {
					if i == 0 {
						tt.given.ordersCreate[i].CustomerID = customer.ID
						order.CustomerID = customer.ID
					}
					err = storage.OrderUpload(ctx, order)
					assert.NoError(t, err, "given:%+v", tt.given)
				}

				//начислить баллы к заказу
				var expectedCurrent float64
				for _, accrual := range tt.given.accruals {
					err = storage.BalanceAccrual(ctx, accrual)
					assert.NoError(t, err, "given:%+v", tt.given)
					expectedCurrent += accrual.Accrual
				}

				//получить текущий баланс покупателя
				balance, err := storage.BalanceByCustomerID(ctx, customer.ID)
				assert.NoError(t, err, "given:%+v", tt.given)

				checkBalance(t, customer.ID, expectedCurrent, 0, balance, tt.given)

				//выполнить списание
				var expectedWithdrawn float64
				for _, withdrawal := range tt.given.withdrawals {

					withdrawal.withdrawal.CustomerID = customer.ID

					err = storage.BalanceWithdraw(ctx, withdrawal.withdrawal)
					assert.ErrorIsf(t, err, withdrawal.expectedErr, "given:%+v", tt.given)

					if err == nil {
						expectedCurrent = balance.Current - withdrawal.withdrawal.Sum
						expectedWithdrawn = balance.Withdrawn + withdrawal.withdrawal.Sum
					}
				}

				//получить текущий баланс покупателя
				balance, err = storage.BalanceByCustomerID(ctx, customer.ID)
				assert.NoError(t, err, "given:%+v", tt.given)

				checkBalance(t, customer.ID, expectedCurrent, expectedWithdrawn, balance, tt.given)

				return errForTransactionRollback
			}, pgx.TxOptions{IsoLevel: pgx.Serializable})

			assert.Error(t, err)
		})
	}
}

func checkBalance(t *testing.T, expectedCustomerID uint64, expectedCurrent float64, expectedWithdrawn float64, balance *model.Balance, given any) {
	t.Helper()

	assert.Equalf(t, expectedCurrent, balance.Current, "balance Current, given:%+v", given)
	assert.Equalf(t, expectedCustomerID, balance.CustomerID, "CustomerID, given:%+v", given)
	assert.Equalf(t, expectedWithdrawn, balance.Withdrawn, "balance Withdrawn, given:%+v", given)
	assert.WithinDurationf(t, time.Now(), balance.UpdatedAt, time.Minute, "balance UpdatedAt, given:%+v", given)
}
