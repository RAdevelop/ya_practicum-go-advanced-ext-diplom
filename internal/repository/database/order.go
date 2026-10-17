package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	"github.com/jackc/pgx/v5"
)

/*
OrderUpload - создание нового заказа

Errors:
  - perror.ErrOrderInvalidModel
  - perror.ErrOrderAlreadyUploadedByCustomer
  - perror.ErrOrderAlreadyUploadedByOther
  - perror.ErrOrderNotFound
  - db error
*/
func (s *Storage) OrderUpload(ctx context.Context, order model.Order) error {
	if !order.IsCorrect() {
		return perror.ErrOrderInvalidModel
	}

	return s.DB.RunInTransaction(ctx, func(ctx context.Context) error {
		// возможно, заказ уже загружен
		orderFound, err := s.orderFindByNumber(ctx, order.Number)
		switch {
		case err == nil:
			if orderFound.CustomerID == order.CustomerID {
				return fmt.Errorf("orderNumber %s, %w", order.Number, perror.ErrOrderAlreadyUploadedByCustomer)
			}
			return fmt.Errorf("orderNumber %s, %w", order.Number, perror.ErrOrderAlreadyUploadedByOther)
		case !errors.Is(err, perror.ErrOrderNotFound):
			return fmt.Errorf("orderNumber %s, %w", order.Number, err)
		}

		// заказа нет — создаём
		err = s.orderCreate(ctx, order)
		if err == nil {
			if !order.Status.IsFinal() {
				// добавляем в очередь
				return s.orderAddToAccrualQueue(ctx, order.Number)
			}
			return nil
		}

		// гонка: кто-то вставил между SELECT и INSERT
		if !isPgErrorCode(err, pgErrUniqueViolationCode) {
			return err
		}

		orderFound, err = s.orderFindByNumber(ctx, order.Number)
		if err != nil {
			return err
		}

		if orderFound.CustomerID == order.CustomerID {
			return fmt.Errorf("orderNumber %s, %w", order.Number, perror.ErrOrderAlreadyUploadedByCustomer)
		}
		return fmt.Errorf("orderNumber %s, %w", order.Number, perror.ErrOrderAlreadyUploadedByOther)
	}, pgx.TxOptions{IsoLevel: pgx.Serializable})
}

/*
OrdersByCustomerID - получение списка заказов указанного покупателя

Errors:
  - perror.ErrCustomerInvalidCredentials
  - perror.ErrOrderNotFound
  - db error
*/
func (s *Storage) OrdersByCustomerID(ctx context.Context, customerID uint64) ([]model.Order, error) {

	if customerID == 0 {
		return nil, fmt.Errorf(`customerID cannot be 0, %w`, perror.ErrCustomerInvalidCredentials)
	}

	/*
		В задаче не было обозначено, как много может быть заказа у пользователя.
		Поэтому
			- не стал делать ограничения на limit/offset по параметрам метода.
			- поставил жесткие OFFSET/LIMIT
	*/

	const sql = `
		SELECT "number", customer_id, status, accrual, uploaded_at, updated_at
		FROM orders
		WHERE customer_id = $1
		ORDER BY uploaded_at DESC
		OFFSET 0 LIMIT 100 
		`

	rows, err := s.DB.Executor(ctx).Query(ctx, sql, customerID)
	if err != nil {
		return nil, fmt.Errorf(`customerID: %d, %w`, customerID, perror.ErrCustomerInvalidCredentials)
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		return nil, fmt.Errorf(`customerID: %d, %w`, customerID, perror.ErrCustomerInvalidCredentials)
	}

	if len(orders) == 0 {
		orders = nil
		return nil, fmt.Errorf(`customerID: %d, %w`, customerID, perror.ErrOrderNotFound)
	}

	return orders, nil
}

/*
orderCreate - создание нового заказа

Errors:
  - db error
*/
func (s *Storage) orderCreate(ctx context.Context, order model.Order) error {

	const sql = `INSERT INTO orders ("number", customer_id, status) VALUES ($1, $2, $3)`
	_, err := s.DB.Executor(ctx).Exec(ctx, sql, order.Number, order.CustomerID, order.Status.String())

	return err
}

/*
orderFindByNumber - поиск заказа по его номеру

Errors:
  - perror.ErrOrderNotFound
  - db error
*/
func (s *Storage) orderFindByNumber(ctx context.Context, orderNumber string) (*model.Order, error) {

	const sql = `
		SELECT "number", customer_id, status, accrual, uploaded_at, updated_at
		FROM orders
		WHERE "number" = $1
	`

	rows, err := s.DB.Executor(ctx).Query(ctx, sql, orderNumber)
	if err != nil {
		return nil, fmt.Errorf("orderNumber: %s, %w", orderNumber, err)
	}

	order, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: orderNumber: %s", perror.ErrOrderNotFound, orderNumber)
		}

		return nil, fmt.Errorf("orderNumber: %s: %w", orderNumber, err)
	}

	return &order, nil
}

/*
orderAccrualUpdate - сохранить начисленные баллы по заказу

Errors:
  - perror.ErrOrderAccrualAlreadyProcessed
  - db error
*/
func (s *Storage) orderAccrualUpdate(ctx context.Context, accrual model.Accrual) error {

	const sql = `
		UPDATE orders
		SET
			status     = $1,
			accrual    = $2,
			updated_at = NOW()
		WHERE "number" = $3 AND status = ANY($4)
	`

	statuses := []string{
		statusInner.AccrualNew.String(),
		statusInner.AccrualProcessing.String(),
	}

	tag, err := s.DB.Executor(ctx).Exec(ctx, sql, accrual.Status.String(), accrual.Accrual, accrual.Order, statuses)
	if err != nil {
		return fmt.Errorf("orderNumber: %s: %w", accrual.Order, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("orderNumber: %s: %w", accrual.Order, perror.ErrOrderAccrualAlreadyProcessed)
	}

	return nil
}
