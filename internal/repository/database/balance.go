package database

import (
	"context"
	"fmt"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/retryer"
	"github.com/jackc/pgx/v5"
)

/*
TODO учитывать:
 - помнить о транзакциях и/или блокировка записей перед начислением/списанием баллов
 - сделать пополнение баланса?! так как при начислении баллов за заказ от внешней системы accrual, эти балы:
  - надо положить в заказ (обновить его со значением accrual)
  - надо начислить эти баллы в текущий баланс покупателя
*/

// BalanceByCustomerID - Получение текущего баланса пользователя
func (s *Storage) BalanceByCustomerID(ctx context.Context, customerID uint64) (*model.Balance, error) {
	//TODO implement
	/*
		#### **Получение текущего баланса пользователя**

		Хендлер: `GET /api/user/balance`.

		Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов лояльности, а также сумме использованных за весь период регистрации баллов.

		Формат запроса:

		```
		GET /api/user/balance HTTP/1.1
		Content-Length: 0
		```

		Возможные коды ответа:

		- `200` — успешная обработка запроса.

		  Формат ответа:

		    ```
		    200 OK HTTP/1.1
		    Content-Type: application/json
		    ...

		    {
		    	"current": 500.5,
		    	"withdrawn": 42
		    }
		    ```

		- `401` — пользователь не авторизован.
		- `500` — внутренняя ошибка сервера.
	*/
	return nil, nil
}

// BalanceWithdrawalsByCustomerID - Получение информации о выводе средств
func (s *Storage) BalanceWithdrawalsByCustomerID(ctx context.Context, customerID uint64) ([]model.Withdrawal, error) {
	//TODO implement

	/*
		#### **Получение информации о выводе средств**

		Хендлер: `GET /api/user/withdrawals`.

		Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых новых к самым старым. Формат даты — RFC3339.

		Формат запроса:

		```
		GET /api/user/withdrawals HTTP/1.1
		Content-Length: 0
		```

		Возможные коды ответа:

		- `200` — успешная обработка запроса.

		  Формат ответа:

		    ```
		    200 OK HTTP/1.1
		    Content-Type: application/json
		    ...

		    [
		        {
		            "order": "2377225624",
		            "sum": 500,
		            "processed_at": "2020-12-09T16:09:57+03:00"
		        }
		    ]
		    ```

		- `204` - нет ни одного списания.
		- `401` — пользователь не авторизован.
		- `500` — внутренняя ошибка сервера.
	*/
	return nil, nil
}

// BalanceWithdraw - списание средств с баланса покупателя
func (s *Storage) BalanceWithdraw(ctx context.Context, withdrawal model.Withdrawal) error {
	//TODO implement
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

// BalanceAccrual - начисление баллов к заказу
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

// balanceCustomerIncrement - увеличиваем баланс покупателя на указанное количество баллов
func (s *Storage) balanceCustomerIncrement(ctx context.Context, customerID uint64, accrual float64) error {

	const sql = `
		INSERT INTO balances (customer_id, current, withdrawn, created_at, updated_at)
		VALUES ($1, $2, 0, NOW(), NOW())
		ON CONFLICT (customer_id)
		DO UPDATE SET
			"current" = "current" + EXCLUDED.current,
			updated_at = NOW()
	`

	_, err := s.DB.Executor(ctx).Exec(ctx, sql, customerID, accrual)
	if err != nil {
		return fmt.Errorf("%w, %w", perror.ErrBalanceIncrement, err)
	}

	return nil
}
