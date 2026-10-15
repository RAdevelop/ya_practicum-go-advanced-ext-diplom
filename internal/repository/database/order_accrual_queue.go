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

	const sql = `
		INSERT INTO order_accrual_queue (order_number)
		VALUES ($1)
		ON CONFLICT (order_number)
			DO NOTHING`

	_, err := s.DB.Executor(ctx).Exec(ctx, sql, orderNumber)
	return err
}

// OrderAccrualQueueUpdate - обновляем состояние очереди для заказов
func (s *Storage) OrderAccrualQueueUpdate(ctx context.Context, orderForUpdate []string, orderForDelete []string) error {

	if len(orderForUpdate) == 0 && len(orderForDelete) == 0 {
		return nil
	}

	return s.DB.RunInTransaction(ctx, func(ctx context.Context) error {

		err := s.orderQueueUpdate(ctx, orderForUpdate)
		if err != nil {
			return err
		}

		err = s.orderQueueDelete(ctx, orderForDelete)
		if err != nil {
			return err
		}

		return nil

	}, pgx.TxOptions{IsoLevel: pgx.Serializable})
}

// orderQueueUpdate - изменяем состояние заказов в очереди
func (s *Storage) orderQueueUpdate(ctx context.Context, orderNumbers []string) error {

	if len(orderNumbers) == 0 {
		return nil
	}

	const sql = `
		UPDATE order_accrual_queue
		SET
			next_check_at = NOW() + (INTERVAL '5 seconds' * POWER(2, LEAST(attempts, 10))),
			attempts      = attempts + 1,
			updated_at    = NOW()
		WHERE order_number = ANY($1)
	`

	_, err := s.DB.Executor(ctx).Exec(ctx, sql, orderNumbers)
	return err
}

// orderQueueDelete - удаляем из очереди заказы, которые получили финальный статус начисления
func (s *Storage) orderQueueDelete(ctx context.Context, orderNumbers []string) error {
	if len(orderNumbers) == 0 {
		return nil
	}

	const sql = `DELETE FROM order_accrual_queue WHERE order_number = ANY($1)`

	_, err := s.DB.Executor(ctx).Exec(ctx, sql, orderNumbers)
	return err
}
