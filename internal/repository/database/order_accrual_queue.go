package database

import (
	"context"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/jackc/pgx/v5"
)

/*
OrderAccrualQueue - заказы, ожидающие начисления

Errors:
  - db error
*/
func (s *Storage) OrderAccrualQueue(ctx context.Context, limit uint) ([]model.OrderAccrualQueue, error) {

	if limit == 0 {
		return nil, nil
	}

	const sql = `
	SELECT order_number, created_at, updated_at, next_check_at, attempts 
	FROM order_accrual_queue
	WHERE next_check_at <= NOW()
	ORDER BY created_at ASC
	OFFSET 0 LIMIT $1
	`

	rows, err := s.DB.Executor(ctx).Query(ctx, sql, limit)
	if err != nil {
		return nil, err
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.OrderAccrualQueue])
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		orders = nil
		return nil, nil
	}

	return orders, nil
}

/*
orderAddToAccrualQueue - добавляем заказ в очередь

Errors:
  - db error
*/
func (s *Storage) orderAddToAccrualQueue(ctx context.Context, orderNumber string) error {

	//TODO del ? created_at, updated_at, next_check_at, attempts
	const sql = `
		INSERT INTO order_accrual_queue (order_number)
		VALUES ($1)
		ON CONFLICT (order_number)
			DO NOTHING`

	_, err := s.DB.Executor(ctx).Exec(ctx, sql, orderNumber)
	return err
}

// OrderAccrualQueueUpdate - обновляем состояние очереди для заказов
func (s *Storage) OrderAccrualQueueUpdate(ctx context.Context, orderNumbers []string) error {
	//TODO implement
	return nil
}

// OrderAccrualQueueDelete - удаляем из очереди заказы, которые получили финальный статус начисления
func (s *Storage) OrderAccrualQueueDelete(ctx context.Context, orderNumbers []string) error {
	//TODO implement
	return nil
}
