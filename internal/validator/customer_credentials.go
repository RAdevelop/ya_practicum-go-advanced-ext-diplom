package validator

import (
	"strings"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
)

/*
IsValidCustomerCredentials - пока простая проверка логина и пароля.
При этом, можно реализовать расширенные требования, не затрагивая работы хендлеров.
*/
func IsValidCustomerCredentials(credentials dto.CustomerCredentials) error {
	if !isValidCredentialsLogin(credentials.Login) || !isValidCredentialsPassword(credentials.Password) {
		return perror.ErrCustomerInvalidCredentials
	}
	return nil
}
func isValidCredentialsLogin(login string) bool {
	return strings.TrimSpace(login) != ""
}
func isValidCredentialsPassword(password string) bool {
	return password != ""
}
