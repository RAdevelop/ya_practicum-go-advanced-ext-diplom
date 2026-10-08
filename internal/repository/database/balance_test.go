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
			name: "ErrBalanceCustomerNotFound",
			given: given{
				customerID:   1,
				ordersCreate: []model.Order{},
				accruals:     []model.Accrual{},
			},
			want: want{
				err: perror.ErrBalanceCustomerNotFound,
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

					assert.Equalf(t, accrualsSum, balance.Current, "given:%+v", tt.given)
					assert.Equalf(t, tt.given.customerID, balance.CustomerID, "given:%+v", tt.given)
					assert.Equalf(t, float64(0), balance.Withdrawn, "given:%+v", tt.given)
					assert.WithinDurationf(t, time.Now(), balance.UpdatedAt, time.Minute, "given:%+v", tt.given)
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
	//TODO implement
	t.Skip("TODO implement")
}
