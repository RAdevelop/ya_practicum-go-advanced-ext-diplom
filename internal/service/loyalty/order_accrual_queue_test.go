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
	"github.com/stretchr/testify/mock"
)

func Test_OrderAccrualQueue(t *testing.T) {

	type given struct {
		inputStatuses        []statusInner.Accrual
		outOrderAccrualQueue []model.OrderAccrualQueue
		makeStorage          func(t *testing.T, outOrdersModel []model.OrderAccrualQueue, outErr error) LoyaltyStorage
	}

	type want struct {
		orders []dto.OrderAccrualQueue
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
				outOrderAccrualQueue: []model.OrderAccrualQueue{
					{
						OrderNumber: "4532015112830366",
						CreatedAt:   time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
						UpdatedAt:   time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
						NextCheckAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
						Attempts:    1,
					},
				},
				makeStorage: func(t *testing.T, outOrdersModel []model.OrderAccrualQueue, outErr error) LoyaltyStorage {
					storage := NewMockLoyaltyStorage(t)
					storage.EXPECT().OrderAccrualQueue(t.Context(), mock.Anything, uint(100)).Return(outOrdersModel, outErr)
					return storage
				},
			},
			want: want{
				err: nil,
				orders: []dto.OrderAccrualQueue{
					{
						OrderNumber: "4532015112830366",
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
				outOrderAccrualQueue: nil,
				makeStorage: func(t *testing.T, outOrdersModel []model.OrderAccrualQueue, outErr error) LoyaltyStorage {
					storage := NewMockLoyaltyStorage(t)
					storage.EXPECT().OrderAccrualQueue(t.Context(), mock.Anything, uint(100)).Return(outOrdersModel, outErr)
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
			loyaltyManager := NewManager(tt.given.makeStorage(t, tt.given.outOrderAccrualQueue, tt.want.err))

			orders, err := loyaltyManager.OrderAccrualQueue(t.Context(), uint(100))
			assert.ErrorIsf(t, err, tt.want.err, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.orders, orders, "given: %+v", tt.given)
		})
	}
}

func Test_OrderAccrualQueueUpdate(t *testing.T) {

	type given struct {
		inputOrderAccrualQueue []dto.OrderAccrualQueue
		inputOrderForUpdate    []string
		inputOrderForDelete    []string
		makeStorage            func(t *testing.T, inputOrderForUpdate []string, inputOrderForDelete []string, err error) LoyaltyStorage
	}

	makeStorage := func(t *testing.T, inputOrderForUpdate []string, inputOrderForDelete []string, err error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		storage.EXPECT().OrderAccrualQueueUpdate(t.Context(), inputOrderForUpdate, inputOrderForDelete).Return(err)
		return storage
	}

	makeStorageOrderAccrualQueueUpdateNever := func(t *testing.T, inputOrderForUpdate []string, inputOrderForDelete []string, err error) LoyaltyStorage {
		storage := NewMockLoyaltyStorage(t)
		return storage
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
				inputOrderAccrualQueue: []dto.OrderAccrualQueue{
					{
						OrderNumber: "4532015112830366",
						NeedUpdate:  true,
					},
					{
						OrderNumber: "4532015112830368",
						NeedUpdate:  false,
					},
				},
				inputOrderForUpdate: []string{
					"4532015112830366",
				},
				inputOrderForDelete: []string{
					"4532015112830368",
				},
				makeStorage: makeStorage,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "StorageOrderAccrualQueueUpdateNeverCalled",
			given: given{
				inputOrderAccrualQueue: []dto.OrderAccrualQueue{},
				inputOrderForUpdate:    []string{},
				inputOrderForDelete:    []string{},
				makeStorage:            makeStorageOrderAccrualQueueUpdateNever,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "error",
			given: given{
				inputOrderAccrualQueue: []dto.OrderAccrualQueue{
					{
						OrderNumber: "4532015112830366",
						NeedUpdate:  true,
					},
					{
						OrderNumber: "4532015112830368",
						NeedUpdate:  false,
					},
				},
				inputOrderForUpdate: []string{
					"4532015112830366",
				},
				inputOrderForDelete: []string{
					"4532015112830368",
				},
				makeStorage: makeStorage,
			},
			want: want{
				err: errors.New("something went wrong"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyManager := NewManager(tt.given.makeStorage(t, tt.given.inputOrderForUpdate, tt.given.inputOrderForDelete, tt.want.err))

			err := loyaltyManager.OrderAccrualQueueUpdate(t.Context(), tt.given.inputOrderAccrualQueue)
			assert.ErrorIsf(t, err, tt.want.err, "given: %+v", tt.given)
		})
	}
}
