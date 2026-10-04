// Package dto - описание структур для запросов
package dto

import (
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

// BalanceWithdraw - для запроса на списание средств
type BalanceWithdraw struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
}

// UserCredentials - для запроса регистрации, авторизации
type UserCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// String - скрыть пароль на случай где-либо его вывести как строку через структуру
func (uc UserCredentials) String() string {
	return fmt.Sprintf("Credentials{Login: %q, Password: ***}", uc.Login)
}

/*
UserJWT - DTO для работы с JWT токеном

Для JSON полей префикс user_* на случай расширения полей, которые можно будет
хранить в токене (чтобы было четкое отличие какое поле к чему относится)
*/
type UserJWT struct {
	UserLogin string `json:"user_login,omitempty"`
	jwt.RegisteredClaims
}

// UserID - парсит Subject в userID
func (jwt *UserJWT) UserID() (uint64, error) {
	userID, err := strconv.ParseUint(jwt.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid subject %q: %w", jwt.Subject, err)
	}
	return userID, nil
}

// User - DTO для передачи данных пользователя между слоями приложения
type User struct {
	ID    uint64 `json:"id"`
	Login string `json:"login"`
}
