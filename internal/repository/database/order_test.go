package database

import (
	"context"
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	"github.com/stretchr/testify/assert"
)

func TestStorage_OrderUpload(t *testing.T) {

	type given struct {
		order model.Order
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
				order: model.Order{
					Number:     "4532015112830366",
					CustomerID: uint64(1),
					Status:     statusInner.AccrualNew,
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "ErrInvalidOrder empty order",
			given: given{},
			want: want{
				err: perror.ErrOrderInvalidModel,
			},
		},
		{
			name: "ErrInvalidOrder empty Number",
			given: given{
				order: model.Order{
					Number:     "",
					CustomerID: uint64(1),
					Status:     statusInner.AccrualNew,
				},
			},
			want: want{
				err: perror.ErrOrderInvalidModel,
			},
		},
		{
			name: "ErrInvalidOrder empty CustomerID",
			given: given{
				order: model.Order{
					Number:     "4532015112830366",
					CustomerID: uint64(0),
					Status:     statusInner.AccrualNew,
				},
			},
			want: want{
				err: perror.ErrOrderInvalidModel,
			},
		},
		{
			name: "ErrInvalidOrder invalid Status",
			given: given{
				order: model.Order{
					Number:     "4532015112830366",
					CustomerID: uint64(1),
					Status:     statusInner.Accrual{},
				},
			},
			want: want{
				err: perror.ErrOrderInvalidModel,
			},
		},
	}

	storage := setUpStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := storage.DB.RunInTransaction(t.Context(), func(ctx context.Context) error {

				err := storage.OrderUpload(ctx, tt.given.order)
				assert.Equalf(t, tt.want.err, err, "given: %+v", tt.given)

				return errForTransactionRollback
			})

			assert.Error(t, err)
		})
	}
}

func TestStorage_Orders(t *testing.T) {

	type given struct {
		customerID uint64
		order      *model.Order
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
				customerID: uint64(1),
				order: &model.Order{
					Number:     "4532015112830366",
					CustomerID: uint64(1),
					Status:     statusInner.AccrualNew,
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "ErrOrderNotFound for customer",
			given: given{
				customerID: uint64(2), // другой покупатель
				order: &model.Order{
					Number:     "4532015112830366",
					CustomerID: uint64(1), // а для этого покупателя есть заказы
					Status:     statusInner.AccrualNew,
				},
			},
			want: want{
				err: perror.ErrOrderNotFound,
			},
		},
		{
			name: "ErrCustomerInvalidCredentials",
			given: given{
				customerID: uint64(0),
			},
			want: want{
				err: perror.ErrCustomerInvalidCredentials,
			},
		},
	}

	storage := setUpStorage(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.DB.RunInTransaction(t.Context(), func(ctx context.Context) error {

				if tt.given.order != nil {
					err := storage.OrderUpload(ctx, *tt.given.order)
					assert.NoErrorf(t, err, "given: %+v", tt.given)
				}

				orders, err := storage.OrdersByCustomerID(ctx, tt.given.customerID)

				assert.ErrorIsf(t, err, tt.want.err, "given: %+v", tt.given)

				if tt.want.err != nil {
					assert.Nilf(t, orders, "given: %+v", tt.given)
				} else {
					assert.Greaterf(t, len(orders), 0, "given: %+v", tt.given)
				}

				return errForTransactionRollback
			})

			assert.Error(t, err)
		})
	}
}

func TestStorage_OrdersAwaitingAccrual(t *testing.T) {
	type given struct {
		inputStatuses []statusInner.Accrual
		inputOrders   []model.Order
	}

	type want struct {
		orderNumbers []string
		err          error
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success find 2 orders",
			given: given{
				inputStatuses: []statusInner.Accrual{
					statusInner.AccrualNew,
					statusInner.AccrualProcessing,
				},
				inputOrders: []model.Order{
					{
						Number:     "4532015112830366",
						CustomerID: uint64(1),
						Status:     statusInner.AccrualNew,
					},
					{
						Number:     "4532015112830367",
						CustomerID: uint64(1),
						Status:     statusInner.AccrualProcessing,
					},
				},
			},
			want: want{
				orderNumbers: []string{
					"4532015112830366",
					"4532015112830367",
				},
				err: nil,
			},
		},
		{
			name: "success find 1 orders",
			given: given{
				inputStatuses: []statusInner.Accrual{
					statusInner.AccrualNew,
					statusInner.AccrualProcessing,
				},
				inputOrders: []model.Order{
					{
						Number:     "4532015112830366",
						CustomerID: uint64(1),
						Status:     statusInner.AccrualNew,
					},
					{
						Number:     "4532015112830367",
						CustomerID: uint64(1),
						Status:     statusInner.AccrualProcessed, //этот заказ не будет найден
					},
				},
			},
			want: want{
				orderNumbers: []string{
					"4532015112830366",
				},
				err: nil,
			},
		},
		{
			name: "ErrOrderNotFound",
			given: given{
				inputStatuses: []statusInner.Accrual{
					statusInner.AccrualNew,
					statusInner.AccrualProcessing,
				},
				inputOrders: []model.Order{
					{
						Number:     "4532015112830366",
						CustomerID: uint64(1),
						Status:     statusInner.AccrualInvalid, //этот заказ не будет найден
						Accrual:    30,
					},
					{
						Number:     "4532015112830367",
						CustomerID: uint64(1),
						Status:     statusInner.AccrualProcessed, //этот заказ не будет найден
						Accrual:    20,
					},
				},
			},
			want: want{
				orderNumbers: []string{},
				err:          perror.ErrOrderNotFound,
			},
		},
	}
	storage := setUpStorage(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := storage.DB.RunInTransaction(t.Context(), func(ctx context.Context) error {

				//добавить заказ(ы)
				for _, inputOrder := range tt.given.inputOrders {
					err := storage.OrderUpload(ctx, inputOrder)
					assert.NoErrorf(t, err, "given: %+v", tt.given)
				}

				//найти заказы в статусе
				orders, err := storage.OrdersAwaitingAccrual(ctx, tt.given.inputStatuses)
				assert.ErrorIsf(t, err, tt.want.err, "given: %+v", tt.given)

				if tt.want.err != nil {
					assert.Nilf(t, orders, "given: %+v", tt.given)
				} else {

					assert.Equalf(t, len(tt.want.orderNumbers), len(orders), "given: %+v", tt.given)
					countEquals := 0
					countEqualsExpected := len(tt.want.orderNumbers)
					for _, order := range orders {
						for _, orderNumber := range tt.want.orderNumbers {
							if orderNumber == order.Number {
								countEquals++
							}
						}
					}
					assert.Equalf(t, countEqualsExpected, countEquals, "given: %+v", tt.given)
				}

				return errForTransactionRollback
			})

			assert.Error(t, err)
		})
	}
}
