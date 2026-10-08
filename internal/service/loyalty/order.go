package loyalty

import (
	"context"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
)

// OrderUpload - загрузка заказа
func (lm *Manager) OrderUpload(ctx context.Context, customerDTO *dto.Customer, number string) error {

	order := model.Order{
		CustomerID: customerDTO.ID,
		Number:     number,
		Status:     statusInner.AccrualNew,
	}
	return lm.storage.OrderUpload(ctx, order)
}

// Orders - получение списка загруженных номеров заказов
func (lm *Manager) Orders(ctx context.Context, customerDTO *dto.Customer) ([]dto.Order, error) {
	ordersModel, err := lm.storage.OrdersByCustomerID(ctx, customerDTO.ID)
	if err != nil {
		return nil, err
	}

	orders := make([]dto.Order, 0, len(ordersModel))
	for _, order := range ordersModel {
		orders = append(orders, dto.Order{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt,
		})
	}
	ordersModel = nil
	return orders, nil
}

// OrdersAwaitingAccrual - заказы, ожидающие начисления
func (lm *Manager) OrdersAwaitingAccrual(statuses []statusInner.Accrual) ([]dto.Order, error) {

	ordersModel, err := lm.storage.OrdersAwaitingAccrual(statuses)
	if err != nil {
		return nil, err
	}

	orders := make([]dto.Order, 0, len(ordersModel))
	for _, order := range ordersModel {
		orders = append(orders, dto.Order{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt,
		})
	}

	ordersModel = nil
	return orders, nil
}
