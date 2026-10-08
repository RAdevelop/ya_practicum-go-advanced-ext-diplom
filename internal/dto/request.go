// Package dto - описание структур для запросов
package dto

import (
	"fmt"
	"strconv"

	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	statusOuter "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/outer"
	"github.com/golang-jwt/jwt/v5"
)

// BalanceWithdraw - для запроса на списание средств
type BalanceWithdraw struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
}

// CustomerCredentials - для запроса регистрации, авторизации
type CustomerCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// String - скрыть пароль на случай где-либо его вывести как строку через структуру
func (cc CustomerCredentials) String() string {
	return fmt.Sprintf("Credentials{Login: %q, Password: ***}", cc.Login)
}

/*
CustomerJWT - DTO для работы с JWT токеном

Для JSON полей префикс customer_* на случай расширения полей, которые можно будет
хранить в токене (чтобы было четкое отличие какое поле к чему относится)
*/
type CustomerJWT struct {
	CustomerLogin string `json:"customer_login,omitempty"`
	jwt.RegisteredClaims
}

// CustomerID - парсит Subject в customerID
func (jwt *CustomerJWT) CustomerID() (uint64, error) {
	customerID, err := strconv.ParseUint(jwt.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid subject %q: %w", jwt.Subject, err)
	}
	return customerID, nil
}

// Customer - DTO для передачи данных пользователя между слоями приложения
type Customer struct {
	ID    uint64 `json:"id"`
	Login string `json:"login"`
}

// AccrualOuter - Начисление от внешней системы
type AccrualOuter struct {
	Order   string              `json:"order"`             // Номер заказа
	Status  statusOuter.Accrual `json:"status"`            // Статус расчёта начисления
	Accrual float64             `json:"accrual,omitempty"` // Рассчитанные баллы к начислению
}

// AccrualInner - Начисление от внешней системы со внутренним статусом
type AccrualInner struct {
	Order   string              `json:"order"`             // Номер заказа
	Status  statusInner.Accrual `json:"status"`            // Статус расчёта начисления
	Accrual float64             `json:"accrual,omitempty"` // Рассчитанные баллы к начислению
}
