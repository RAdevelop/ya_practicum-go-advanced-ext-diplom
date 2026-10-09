package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/retryer"
	"github.com/jackc/pgx/v5"
)

/*
TODO учитывать:
 - помнить о транзакциях и/или блокировка записей перед начислением/списанием баллов
*/

/*
BalanceByCustomerID - Получение текущего баланса пользователя

Errors:
  - perror.ErrCustomerInvalidCredentials
  - perror.ErrBalanceCustomerNotFound
  - db error
*/
func (s *Storage) BalanceByCustomerID(ctx context.Context, customerID uint64) (*model.Balance, error) {

	if customerID == 0 {
		return nil, perror.ErrCustomerInvalidCredentials
	}

	const sql = `
		SELECT id, customer_id, "current", withdrawn, updated_at
		FROM balances
		WHERE customer_id = $1
	`

	rows, err := s.DB.Executor(ctx).Query(ctx, sql, customerID)
	if err != nil {
		return nil, err
	}

	balance, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Balance])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w, %w", perror.ErrBalanceCustomerNotFound, err)
		}
		return nil, err
	}

	return &balance, nil
}

/*
BalanceWithdrawalsByCustomerID - Получение информации о выводе средств

Errors:
  - perror.ErrCustomerInvalidCredentials
  - perror.ErrWithdrawalNotFound
  - db error
*/
func (s *Storage) BalanceWithdrawalsByCustomerID(ctx context.Context, customerID uint64) ([]model.Withdrawal, error) {

	if customerID == 0 {
		return nil, perror.ErrCustomerInvalidCredentials
	}

	const sql = `
		SELECT w.id, w.order_id, w.sum, w.processed_at, o.number AS "order", o.customer_id
		FROM withdrawals AS w 
		JOIN orders AS o ON(o.id = w.order_id)
		WHERE o.customer_id = $1
		ORDER BY w.processed_at DESC
	`

	rows, err := s.DB.Executor(ctx).Query(ctx, sql, customerID)
	if err != nil {
		return nil, err
	}

	withdrawals, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Withdrawal])
	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, perror.ErrWithdrawalNotFound
	}

	return withdrawals, nil
}

/*
BalanceWithdraw - списание средств с баланса покупателя

Errors:
- perror.
*/
func (s *Storage) BalanceWithdraw(ctx context.Context, withdrawal model.Withdrawal) error {
	//TODO implement - помнить о транзакциях и/или блокировка записей перед начислением/списанием баллов
	/*
			#### **Запрос на списание средств**

			Хендлер: `POST /api/user/balance/withdraw`

			Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя в счет оплаты которого списываются баллы.

			Примечание: для успешного списания достаточно успешной регистрации запроса, никаких внешних систем начисления не предусмотрено и не требуется реализовывать.

			Формат запроса:

			```
			POST /api/user/balance/withdraw HTTP/1.1
			Content-Type: application/json

			{
				"order": "2377225624",
			    "sum": 751
			}
			```

			Здесь `order` — номер заказа, а `sum` — сумма баллов к списанию в счёт оплаты.

			Возможные коды ответа:

			- `200` — успешная обработка запроса;
			- `401` — пользователь не авторизован;
			- `402` — на счету недостаточно средств;
			- `422` — неверный номер заказа;
			- `500` — внутренняя ошибка сервера.

		TODO надо будет
		 - перерассчитывать текущий баланс покупателя при успешном списании баллов
		 - помнить о транзакции и/или блокировки записей в таблицАХ!
		 - помнить, что баланс не должен уходить в минус. Если расчет показывает отрицательное значение - возвращать ошибку (не достаточно средств)!
	*/
	return nil
}

/*
BalanceAccrual - начисление баллов к заказу

ВАЖНО: использует внутри себя транзакцию с повторными попытками!

Errors:
  - perror.ErrOrderAccrualAlreadyProcessed
  - perror.ErrOrderNotFound
  - perror.ErrBalanceIncrement
  - db error
*/
func (s *Storage) BalanceAccrual(ctx context.Context, accrual model.Accrual) error {
	/*
		- найти заказ по номеру
			- если не нашли, выходим с ошибкой
			- если нашли, то знаем id покупателя
		- обновить состояние заказа по номеру заказа
			- если не смогли, выходим с ошибкой
		- обновить баланс покупателя по id покупателя
			- если не смогли, выходим с ошибкой
	*/

	// на случай dead-lock или аналогичных ситуаций
	_, err := retryer.RetryLinear(ctx, func(ctx context.Context) (struct{}, error) {
		err := s.DB.RunInTransaction(ctx, func(ctx context.Context) error {

			order, err := s.orderFindByNumber(ctx, accrual.Order)
			if err != nil {
				return err
			}

			if order.Status.IsFinal() {
				return perror.ErrOrderAccrualAlreadyProcessed
			}

			err = s.orderAccrualUpdate(ctx, accrual)
			if err != nil {
				return err
			}

			err = s.balanceCustomerIncrement(ctx, order.CustomerID, accrual.Accrual)
			if err != nil {
				return err
			}

			return nil
		}, pgx.TxOptions{IsoLevel: pgx.Serializable})

		return struct{}{}, err
	}, uint(1), new(uint(3)))

	return err
}

/*
balanceCustomerIncrement - увеличиваем баланс покупателя на указанное количество баллов

Errors:
  - perror.ErrBalanceIncrement
  - db error
*/
func (s *Storage) balanceCustomerIncrement(ctx context.Context, customerID uint64, accrual float64) error {

	const sql = `
		INSERT INTO balances (customer_id, "current", withdrawn, updated_at)
		VALUES ($1, $2, 0, NOW())
		ON CONFLICT (customer_id)
		DO UPDATE SET
			"current" = balances."current" + EXCLUDED."current",
			updated_at = NOW()
	`

	_, err := s.DB.Executor(ctx).Exec(ctx, sql, customerID, accrual)
	if err != nil {
		return fmt.Errorf("%w, %w", perror.ErrBalanceIncrement, err)
	}

	return nil
}
