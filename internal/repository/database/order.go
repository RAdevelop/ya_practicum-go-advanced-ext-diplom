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
TODO судя по всему, надо будет делать методы еще для:
 - получения списка заказов по статусам (для их отправки во внешнюю систему расчетов)
 - обновления данных по заказу (после того, как будут получена информация о начислении баллов за закза)
 - имеет смысл делать отдельную таблицу для начисления баллов. Так как есть (будет) таблица для истории списания баллов.
  - тогда поле "accrual" из таблицы orders уйдет в такую таблицу истории начисления.
  - тогда для возврата списка заказов надо делать join (модель заказа + модель истории начисления баллов?!) - возвращать DTO, а не модель!?
  - тогда надо будет такой "json" собирать на уровне сервиса?!
	{
            "number": "9278923470",
            "status": "PROCESSED",
            "accrual": 500,
            "uploaded_at": "2020-12-10T15:15:45+03:00"
        },
  - в таблице orders поле ID скорее всего избыточное! так как таблица истории заказов будет "соединяться" все равно по "номеру заказа".
   - иначе придется делать сначала поиск заказа по номеру, потом получать его ID, потом вставлять в таблицу истории?!
   - с другой стороны, все равно надо делать "поиск" заказа по его номеру, иначе в истории начисления могут быть "фантомные" записи, не привязанные ни к одному заказу!
    - а в этом случае, ID заказа уже будет известен...
    - такую операцию делать в транзакции?! :
	  - найти заказ
	  - положить балы в историю начисления
	  - обновить текущий баланс пользователя
*/

func (s *Storage) OrderUpload(ctx context.Context, order model.Order) error {
	if !order.IsCorrect() {
		return perror.ErrOrderInvalidModel
	}

	// возможно, заказ уже загружен.
	orderFound, err := s.orderFindByNumber(ctx, order.Number)
	switch {
	case err == nil:
		if orderFound.CustomerID == order.CustomerID {
			return perror.ErrOrderAlreadyUploadedByCustomer
		}
		return perror.ErrOrderAlreadyUploadedByOther
	case !errors.Is(err, perror.ErrOrderNotFound):
		return err
	}

	// заказа нет — создаём.
	err = s.orderCreate(ctx, order)
	if err == nil {
		return nil
	}

	// гонка: кто-то вставил между SELECT и INSERT.
	if !isPgErrorCode(err, pgErrUniqueViolationCode) {
		return err
	}

	orderFound, err = s.orderFindByNumber(ctx, order.Number)
	if err != nil {
		return err
	}

	if orderFound.CustomerID == order.CustomerID {
		return perror.ErrOrderAlreadyUploadedByCustomer
	}
	return perror.ErrOrderAlreadyUploadedByOther
}

func (s *Storage) OrdersByCustomerID(ctx context.Context, customerID uint64) ([]model.Order, error) {

	if customerID == 0 {
		return nil, perror.ErrCustomerInvalidCredentials
	}

	/*
		В задаче не было обозначено, как много может быть заказа у пользователя.
		Поэтому
			- не стал делать ограничения на limit/offset по параметрам метода.
			- поставил жесткие OFFSET/LIMIT
	*/

	const sql = `
		SELECT id, "number", customer_id, status, accrual, uploaded_at, updated_at
		FROM orders
		WHERE customer_id = $1
		ORDER BY uploaded_at DESC
		OFFSET 0 LIMIT 100 
		`

	rows, err := s.DB.Executor(ctx).Query(ctx, sql, customerID)

	if err != nil {
		return nil, err
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		orders = nil
		return nil, fmt.Errorf("%w", perror.ErrOrderNotFound)
	}

	return orders, nil
}

// OrdersAwaitingAccrual - заказы, ожидающие начисления
func (s *Storage) OrdersAwaitingAccrual(ctx context.Context, statuses []statusInner.Accrual) ([]model.Order, error) {

	const sql = `
	SELECT id, "number", customer_id, status, accrual, uploaded_at, updated_at FROM orders  WHERE status = ANY($1)
	ORDER BY uploaded_at DESC
	OFFSET 0 LIMIT 100
	`

	orderStatuses := make([]string, 0, len(statuses))
	for _, status := range statuses {
		orderStatuses = append(orderStatuses, status.String())
	}

	rows, err := s.DB.Executor(ctx).Query(ctx, sql, orderStatuses)
	if err != nil {
		return nil, err
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		orders = nil
		return nil, perror.ErrOrderNotFound
	}

	return orders, nil
}

// orderCreate - создание нового заказа
func (s *Storage) orderCreate(ctx context.Context, order model.Order) error {

	const sql = `INSERT INTO orders ("number", customer_id, status) VALUES ($1, $2, $3)`
	_, err := s.DB.Executor(ctx).Exec(ctx, sql, order.Number, order.CustomerID, order.Status.String())

	return err
}

// orderFindByNumber - поиск заказа по его номеру
func (s *Storage) orderFindByNumber(ctx context.Context, number string) (*model.Order, error) {

	const sql = `
		SELECT id, "number", customer_id, status, accrual, uploaded_at, updated_at
		FROM orders
		WHERE "number" = $1
	`

	rows, err := s.DB.Executor(ctx).Query(ctx, sql, number)
	if err != nil {
		return nil, err
	}

	order, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, perror.ErrOrderNotFound
		}

		return nil, err
	}

	return &order, nil
}

// orderAccrualUpdate - сохранить начисленные баллы по заказу
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
		return fmt.Errorf("%w, %w", perror.ErrAccrualApply, err)
	}

	if tag.RowsAffected() == 0 {
		return perror.ErrOrderAccrualAlreadyProcessed
	}

	return nil
}
