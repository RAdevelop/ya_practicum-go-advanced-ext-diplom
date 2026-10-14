package loyalty

import (
	"context"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
)

// OrderAccrualQueue - заказы, ожидающие начисления
func (lm *Manager) OrderAccrualQueue(ctx context.Context, limit uint) ([]dto.OrderAccrualQueue, error) {

	orderAccrualQueue, err := lm.storage.OrderAccrualQueue(ctx, limit)
	if err != nil {
		return nil, err
	}

	if len(orderAccrualQueue) == 0 {
		return nil, nil
	}

	orderQueue := make([]dto.OrderAccrualQueue, 0, len(orderAccrualQueue))
	for _, order := range orderAccrualQueue {
		orderQueue = append(orderQueue, dto.OrderAccrualQueue{
			OrderNumber: order.OrderNumber,
		})
	}

	orderAccrualQueue = nil
	return orderQueue, nil
}

// OrderAccrualQueueUpdate - обновляем состояние очереди для заказов
func (lm *Manager) OrderAccrualQueueUpdate(ctx context.Context, orderNumbers []string) error {
	return lm.storage.OrderAccrualQueueUpdate(ctx, orderNumbers)
}

// OrderAccrualQueueDelete - удаляем из очереди заказы, которые получили финальный статус начисления
func (lm *Manager) OrderAccrualQueueDelete(ctx context.Context, orderNumbers []string) error {
	return lm.storage.OrderAccrualQueueDelete(ctx, orderNumbers)
}
