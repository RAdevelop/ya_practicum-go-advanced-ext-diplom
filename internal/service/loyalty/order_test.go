package loyalty

import (
	"errors"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	"github.com/stretchr/testify/assert"
)

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
			loyaltyManager := NewManager(tt.given.makeStorage(t, tt.given.inputCustomerDTO, tt.given.inputOrderNumber, tt.want.outErr))

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
		orders []dto.Order
	}

	makeStorage := func(t *testing.T, inputCustomerDTO *dto.Customer, outOrders []model.Order, outErr error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().OrdersByCustomerID(t.Context(), inputCustomerDTO.ID).Return(outOrders, outErr)
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
				outOrders: []model.Order{
					{
						ID:         uint64(1),
						Number:     "4532015112830366",
						CustomerID: uint64(1),
						Accrual:    0,
						Status:     statusInner.AccrualNew,
						UploadedAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
						UpdatedAt:  time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
					},
				},
				outErr:      nil,
				makeStorage: makeStorage,
			},
			want: want{
				orders: []dto.Order{
					{
						Number:     "4532015112830366",
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
			loyaltyManager := NewManager(tt.given.makeStorage(t, tt.given.inputCustomerDTO, tt.given.outOrders, tt.given.outErr))

			orders, err := loyaltyManager.Orders(t.Context(), tt.given.inputCustomerDTO)

			assert.ErrorIsf(t, err, tt.given.outErr, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.orders, orders, "given: %+v", tt.given)
		})
	}
}

func Test_OrdersAwaitingAccrual(t *testing.T) {

	type given struct {
		inputStatuses  []statusInner.Accrual
		outOrdersModel []model.Order
		makeStorage    func(t *testing.T, inputStatuses []statusInner.Accrual, outOrdersModel []model.Order, outErr error) LoyaltyStorage
	}

	type want struct {
		orders []dto.Order
		err    error
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success",
			given: given{
				inputStatuses: []statusInner.Accrual{
					statusInner.AccrualNew,
					statusInner.AccrualProcessing,
				},
				outOrdersModel: []model.Order{
					{
						ID:         uint64(1),
						Number:     "4532015112830366",
						CustomerID: uint64(1),
						Status:     statusInner.AccrualNew,
						Accrual:    float64(0),
						UploadedAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
						UpdatedAt:  time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
					},
				},
				makeStorage: func(t *testing.T, inputStatuses []statusInner.Accrual, outOrdersModel []model.Order, outErr error) LoyaltyStorage {
					storage := NewMockLoyaltyStorage(t)
					storage.EXPECT().OrdersAwaitingAccrual(t.Context(), inputStatuses).Return(outOrdersModel, outErr)
					return storage
				},
			},
			want: want{
				err: nil,
				orders: []dto.Order{
					{
						Number:     "4532015112830366",
						Status:     statusInner.AccrualNew,
						UploadedAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
					},
				},
			},
		},
		{
			name: "error",
			given: given{
				inputStatuses: []statusInner.Accrual{
					statusInner.AccrualNew,
					statusInner.AccrualProcessing,
				},
				outOrdersModel: nil,
				makeStorage: func(t *testing.T, inputStatuses []statusInner.Accrual, outOrdersModel []model.Order, outErr error) LoyaltyStorage {
					storage := NewMockLoyaltyStorage(t)
					storage.EXPECT().OrdersAwaitingAccrual(t.Context(), inputStatuses).Return(outOrdersModel, outErr)
					return storage
				},
			},
			want: want{
				err:    perror.ErrOrderNotFound,
				orders: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewManager(tt.given.makeStorage(t, tt.given.inputStatuses, tt.given.outOrdersModel, tt.want.err))

			orders, err := loyaltyManager.OrdersAwaitingAccrual(t.Context(), tt.given.inputStatuses)
			assert.ErrorIsf(t, err, tt.want.err, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.orders, orders, "given: %+v", tt.given)
		})
	}
}
