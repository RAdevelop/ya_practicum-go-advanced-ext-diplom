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

				orders, err := storage.Orders(ctx, tt.given.customerID)

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
