package database

import (
	"context"
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	"github.com/stretchr/testify/assert"
)

func TestStorage_OrderAccrualQueue(t *testing.T) {
	type given struct {
		inputStatuses []statusInner.Accrual
		inputOrders   []model.Order
	}

	type want struct {
		limit        uint
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
				limit: 2,
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
				limit: 100,
				orderNumbers: []string{
					"4532015112830366",
				},
				err: nil,
			},
		},
		{
			name: "empty list",
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
				limit:        100,
				orderNumbers: []string{},
				err:          nil,
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

				//найти заказы в очереди
				orders, err := storage.OrderAccrualQueue(ctx, tt.want.limit)
				assert.ErrorIsf(t, err, tt.want.err, "given: %+v", tt.given)

				if tt.want.err != nil {
					assert.Nilf(t, orders, "given: %+v", tt.given)
				} else {

					assert.Equalf(t, len(tt.want.orderNumbers), len(orders), "given: %+v", tt.given)
					countEquals := 0
					countEqualsExpected := len(tt.want.orderNumbers)
					for _, order := range orders {
						for _, orderNumber := range tt.want.orderNumbers {
							if orderNumber == order.OrderNumber {
								countEquals++
							}
						}
					}
					assert.Equalf(t, countEqualsExpected, countEquals, "given: %+v", tt.given)
					assert.LessOrEqual(t, uint(countEquals), tt.want.limit, "given: %+v", tt.given)
				}

				return errForTransactionRollback
			})

			assert.Error(t, err)
		})
	}
}
