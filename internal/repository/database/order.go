package database

import (
	"context"
	"fmt"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) OrderUpload(ctx context.Context, order model.Order) error {

	if !order.IsCorrect() {
		return perror.ErrOrderInvalidModel
	}

	sql := `INSERT INTO orders ("number", customer_id, status) VALUES ($1, $2, $3)`
	_, err := s.DB.Executor(ctx).Exec(ctx, sql, order.Number, order.CustomerID, order.Status.String())

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Orders(ctx context.Context, customerID uint64) ([]model.Order, error) {

	if customerID == 0 {
		return nil, perror.ErrCustomerInvalidCredentials
	}

	/*
		В задаче не было обозначено, как много может быть заказа у пользователя.
		Поэтому
			- не стал делать ограничения на limit/offset по параметрам метода.
			- поставил жесткие OFFSET/LIMIT
	*/

	sql := `
		SELECT id, "number", customer_id, status, accrual, uploaded_at, updated_at
		FROM orders
		WHERE customer_id = $1
		ORDER BY uploaded_at DESC
		OFFSET 0 LIMIT 100 
		`

	rows, err := s.DB.Executor(ctx).Query(ctx, sql, customerID)

	if err != nil {
		return nil, fmt.Errorf("%w, %w", perror.ErrOrderNotFound, err)
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		return nil, fmt.Errorf("%w, %w", perror.ErrOrderNotFound, err)
	}

	if len(orders) == 0 {
		orders = nil
		return nil, fmt.Errorf("%w, %w", perror.ErrOrderNotFound, err)
	}

	return orders, nil
}
