package database

import (
	"context"
	"errors"
	"fmt"
	"time"

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
  - db error
*/
func (s *Storage) BalanceByCustomerID(ctx context.Context, customerID uint64) (*model.Balance, error) {

	if customerID == 0 {
		return nil, perror.ErrCustomerInvalidCredentials
	}

	const sql = `
		SELECT customer_id, "current", withdrawn, updated_at
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
			return &model.Balance{CustomerID: customerID, UpdatedAt: time.Now()}, nil
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
		SELECT w.order_id, w.sum, w.processed_at, o.number AS "order", o.customer_id
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
  - perror.ErrOrderNotFound
  - perror.ErrBalanceInsufficient
  - perror.ErrWithdrawalAlreadyProcessed
  - db error
*/
func (s *Storage) BalanceWithdraw(ctx context.Context, withdrawal model.Withdrawal) error {

	if !withdrawal.IsCorrect() {
		return perror.ErrWithdrawalInvalid
	}

	_, err := retryer.RetryLinear(ctx, func(ctx context.Context) (struct{}, error) {
		err := s.DB.RunInTransaction(ctx, func(ctx context.Context) error {

			//проверить, что заказ существует и он принадлежит покупателю
			order, err := s.orderFindByNumber(ctx, withdrawal.Order)
			if err != nil {
				return err
			}
			if order.CustomerID != withdrawal.CustomerID || !order.Status.IsFinal() {
				return perror.ErrOrderNotFound
			}

			//получить текущий баланс покупателя
			balance, err := s.BalanceByCustomerID(ctx, withdrawal.CustomerID)
			if err != nil {
				return err
			}
			//проверить, что значение списания меньше или равно, чем есть на балансе у покупателя
			if withdrawal.Sum > balance.Current {
				return perror.ErrBalanceInsufficient
			}
			//добавить списание в историю
			withdrawal.OrderID = order.ID
			err = s.balanceWithdrawHistoryUpdate(ctx, withdrawal)
			if err != nil {
				return err
			}
			//обновить баланс покупателя
			err = s.balanceCustomerDecrement(ctx, withdrawal)
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
  - perror.ErrCustomerInvalidCredentials
  - perror.ErrBalanceIncrement
  - perror.ErrAccrualApply
  - db error
*/
func (s *Storage) balanceCustomerIncrement(ctx context.Context, customerID uint64, accrual float64) error {

	if customerID == 0 {
		return perror.ErrCustomerInvalidCredentials
	}

	if accrual <= 0 {
		return perror.ErrAccrualApply
	}

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

/*
balanceWithdrawHistoryUpdate - добавление информации о списании по заказу

Errors:
  - perror.ErrWithdrawalAlreadyProcessed
  - db error
*/
func (s *Storage) balanceWithdrawHistoryUpdate(ctx context.Context, withdrawal model.Withdrawal) error {

	if !withdrawal.IsCorrect() {
		return perror.ErrWithdrawalInvalid
	}

	const sql = `
 		INSERT INTO withdrawals(order_id, "sum")
 		VALUES ($1, $2)
	`
	_, err := s.DB.Executor(ctx).Exec(ctx, sql, withdrawal.OrderID, withdrawal.Sum)
	if err != nil {
		if isPgErrorCode(err, pgErrUniqueViolationCode) {
			return perror.ErrWithdrawalAlreadyProcessed
		}
		return err
	}

	return nil
}

/*
balanceCustomerDecrement - списание с баланса покупателя

Errors:
  - perror.ErrWithdrawalInvalid
  - perror.ErrBalanceDecrement
  - db error
*/
func (s *Storage) balanceCustomerDecrement(ctx context.Context, withdrawal model.Withdrawal) error {

	if !withdrawal.IsCorrect() {
		return perror.ErrWithdrawalInvalid
	}

	const sql = `
		UPDATE balances
		SET 
		    "current"   = "current" - $1,
		    "withdrawn" = "withdrawn" + $1,
		    updated_at = NOW()
		WHERE customer_id = $2 AND "current" >= $1
	`

	tag, err := s.DB.Executor(ctx).Exec(ctx, sql, withdrawal.Sum, withdrawal.CustomerID)
	if err != nil {
		return fmt.Errorf("%w, %w", perror.ErrBalanceDecrement, err)
	}

	if tag.RowsAffected() == 0 {
		return perror.ErrBalanceDecrement
	}

	return nil
}
