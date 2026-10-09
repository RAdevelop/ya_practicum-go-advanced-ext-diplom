package loyalty

import (
	"context"
	"fmt"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/jwtoken"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"golang.org/x/crypto/bcrypt"
)

// UserLogin - авторизация покупателя
func (lm *Manager) UserLogin(ctx context.Context, customerCredentials dto.CustomerCredentials, jwtSecret string) (string, error) {

	customerModel, err := lm.storage.CustomerFindByLogin(ctx, customerCredentials.Login)

	if err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrCustomerNotFound, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(customerModel.PasswordHash), []byte(customerCredentials.Password)); err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrCustomerInvalidCredentials, err)
	}

	customerDTO := dto.Customer{
		ID:    customerModel.ID,
		Login: customerModel.Login,
	}
	token, err := jwtoken.Generate(customerDTO, []byte(jwtSecret))

	if err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrTokenInvalid, err)
	}

	return token, nil
}

// UserRegister - регистрация покупателя
func (lm *Manager) UserRegister(ctx context.Context, customerCredentials dto.CustomerCredentials, jwtSecret string) (string, error) {

	// 1. Хешируем пароль.
	hash, err := bcrypt.GenerateFromPassword([]byte(customerCredentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrHashGenerateFromPassword, err)
	}

	// 2. Создаём пользователя.
	customerModel := &model.Customer{
		Login:        customerCredentials.Login,
		PasswordHash: string(hash),
	}

	customerModel, err = lm.storage.CustomerCreate(ctx, customerModel)

	if err != nil {
		return "", err
	}

	/*
		Можно считать, что идет дублирование данных.
		Зато слой сервиса генерации токена ничего не знает именно о модели пользователя.
		И в этом случае необходимый рефакторинг слоев можно будет делать независимо.
	*/
	customerDTO := dto.Customer{
		Login: customerModel.Login,
		ID:    customerModel.ID,
	}

	// 3. выдаём JWT.
	token, err := jwtoken.Generate(customerDTO, []byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrTokenInvalid, err)
	}

	return token, nil
}
