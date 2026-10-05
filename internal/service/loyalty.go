package service

/*
Сергей, приветствую!
Да, я отправил ссылку на проверку первый раз без полной реализации всего. Так как для спринта 5 была рекомендация во введении:
"Первую ссылку рекомендуем отправить уже в первые две недели. Необязательно отправлять всю работу сразу. При первой отправке ваша задача — получить корректирующий фидбэк."
:)
Забыл такой комментарий оставить сразу.
Жаль, что при отправке ссылок на задание, нет возможности написать что-то вроде сопроводительного сообщения в GUI Я.Практикума.

И да, не вся реализация есть, так как решил все же пройти по классическому пути TDD.
Поэтому, по твоей рекомендации из прошлых спринтов, сразу стал использовать Мокери (как ты и заметил уже).
По той же причине есть у меня в коде TODO, которые, конечно же буду закрывать. Пишу их для себя.
По PR - "обновил" мастер ветку, привел ее в исходное состояние. До этого поторопился пушить в нее.
В общем, теперь должны быть видны все изменения.
*/

import (
	"context"
	"fmt"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/jwtoken"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"golang.org/x/crypto/bcrypt"
)

// LoyaltyStorage - интерфейс для определения работы с хранилищем данных по программе лояльности
//
//go:generate mockery
type LoyaltyStorage interface {
	OrderUpload(ctx context.Context, customerID uint64, number string) error
	Orders(ctx context.Context, customerID uint64) ([]model.Order, error)
	Balance(ctx context.Context, customerID uint64) (*model.Balance, error)
	BalanceWithdrawals(ctx context.Context, customerID uint64) ([]model.Withdrawal, error)
	BalanceWithdraw(ctx context.Context, customerID uint64, orderNumber string, sum float64) error
	CustomerCreate(ctx context.Context, customer *model.Customer) (*model.Customer, error)
	CustomerFindByLogin(ctx context.Context, login string) (*model.Customer, error)
}

// LoyaltyManager - сервис для работы с программой лояльности
type LoyaltyManager struct {
	storage LoyaltyStorage
}

func NewLoyaltyManager(storage LoyaltyStorage) *LoyaltyManager {
	return &LoyaltyManager{
		storage: storage,
	}
}

// OrderUpload - Загрузка заказа
func (lm *LoyaltyManager) OrderUpload(ctx context.Context, customerDTO *dto.Customer, number string) error {
	return lm.storage.OrderUpload(ctx, customerDTO.ID, number)
}

// Orders - Получение списка загруженных номеров заказов
func (lm *LoyaltyManager) Orders(ctx context.Context, customerDTO *dto.Customer) ([]model.Order, error) {
	return lm.storage.Orders(ctx, customerDTO.ID)
}

// Balance - Получение текущего баланса пользователя
func (lm *LoyaltyManager) Balance(ctx context.Context, customerDTO *dto.Customer) (*model.Balance, error) {
	return lm.storage.Balance(ctx, customerDTO.ID)
}

// BalanceWithdrawals - Получение информации о выводе средств
func (lm *LoyaltyManager) BalanceWithdrawals(ctx context.Context, customerDTO *dto.Customer) ([]model.Withdrawal, error) {
	return lm.storage.BalanceWithdrawals(ctx, customerDTO.ID)
}

// BalanceWithdraw - списание средств
func (lm *LoyaltyManager) BalanceWithdraw(ctx context.Context, customerDTO *dto.Customer, orderNumber string, sum float64) error {
	return lm.storage.BalanceWithdraw(ctx, customerDTO.ID, orderNumber, sum)
}

func (lm *LoyaltyManager) UserLogin(ctx context.Context, customerCredentials dto.CustomerCredentials, jwtSecret string) (string, error) {

	customerModel, err := lm.storage.CustomerFindByLogin(ctx, customerCredentials.Login)

	if err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrCustomerNotFound, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(customerModel.PasswordHash), []byte(customerCredentials.Password)); err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrInvalidCustomerCredentials, err)
	}

	customerDTO := dto.Customer{
		ID:    customerModel.ID,
		Login: customerModel.Login,
	}
	token, err := jwtoken.Generate(customerDTO, []byte(jwtSecret))

	if err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrInvalidToken, err)
	}

	return token, nil
}

func (lm *LoyaltyManager) UserRegister(ctx context.Context, customerCredentials dto.CustomerCredentials, jwtSecret string) (string, error) {

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
		return "", fmt.Errorf("%w: %w", perror.ErrInvalidToken, err)
	}

	return token, nil
}
