package validator

import (
	"strings"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
)

/*
IsValidUserCredentials - пока простая проверка логина и пароля.
При этом, можно реализовать расширенные требования, не затрагивая работы хендлеров.
*/
func IsValidUserCredentials(credentials dto.UserCredentials) error {
	if !isValidUserCredentialsLogin(credentials.Login) || !isValidUserCredentialsPassword(credentials.Password) {
		return perror.ErrInvalidUserCredentials
	}
	return nil
}
func isValidUserCredentialsLogin(login string) bool {
	return strings.Trim(login, " ") != ""
}
func isValidUserCredentialsPassword(password string) bool {
	return password != ""
}
