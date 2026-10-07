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
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	"golang.org/x/crypto/bcrypt"
)

// LoyaltyStorage - интерфейс для определения работы с хранилищем данных по программе лояльности
//
//go:generate mockery
type LoyaltyStorage interface {
	OrderUpload(ctx context.Context, order model.Order) error
	OrdersByCustomerID(ctx context.Context, customerID uint64) ([]model.Order, error)
	BalanceByCustomerID(ctx context.Context, customerID uint64) (*model.Balance, error)
	BalanceWithdrawalsByCustomerID(ctx context.Context, customerID uint64) ([]model.Withdrawal, error)
	BalanceWithdraw(ctx context.Context, withdrawal model.Withdrawal) error
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

	order := model.Order{
		CustomerID: customerDTO.ID,
		Number:     number,
		Status:     statusInner.AccrualNew,
	}
	return lm.storage.OrderUpload(ctx, order)
}

// Orders - Получение списка загруженных номеров заказов
func (lm *LoyaltyManager) Orders(ctx context.Context, customerDTO *dto.Customer) ([]dto.Order, error) {
	ordersModel, err := lm.storage.OrdersByCustomerID(ctx, customerDTO.ID)
	if err != nil {
		return nil, err
	}

	orders := make([]dto.Order, 0, len(ordersModel))
	for _, order := range ordersModel {
		orders = append(orders, dto.Order{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt,
		})
	}
	ordersModel = nil
	return orders, nil
}

// Balance - Получение текущего баланса пользователя
func (lm *LoyaltyManager) Balance(ctx context.Context, customerDTO *dto.Customer) (dto.Balance, error) {
	balanceModel, err := lm.storage.BalanceByCustomerID(ctx, customerDTO.ID)
	if err != nil {
		return dto.Balance{}, err
	}

	balance := dto.Balance{
		Current:   balanceModel.Current,
		Withdrawn: balanceModel.Withdrawn,
	}

	balanceModel = nil
	return balance, nil
}

// BalanceWithdrawals - Получение информации о выводе средств
func (lm *LoyaltyManager) BalanceWithdrawals(ctx context.Context, customerDTO *dto.Customer) ([]dto.Withdrawal, error) {

	withdrawalsModel, err := lm.storage.BalanceWithdrawalsByCustomerID(ctx, customerDTO.ID)
	if err != nil {
		return nil, err
	}

	withdrawals := make([]dto.Withdrawal, 0, len(withdrawalsModel))
	for _, withdrawal := range withdrawalsModel {
		withdrawals = append(withdrawals, dto.Withdrawal{
			Order:       withdrawal.Order,
			Sum:         withdrawal.Sum,
			ProcessedAt: withdrawal.ProcessedAt,
		})
	}

	withdrawalsModel = nil
	return withdrawals, nil
}

// BalanceWithdraw - списание средств
func (lm *LoyaltyManager) BalanceWithdraw(ctx context.Context, customerDTO *dto.Customer, orderNumber string, sum float64) error {

	withdrawal := model.Withdrawal{
		CustomerID: customerDTO.ID,
		Order:      orderNumber,
		Sum:        sum,
	}

	return lm.storage.BalanceWithdraw(ctx, withdrawal)
}

func (lm *LoyaltyManager) UserLogin(ctx context.Context, customerCredentials dto.CustomerCredentials, jwtSecret string) (string, error) {

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
		return "", fmt.Errorf("%w: %w", perror.ErrTokenInvalid, err)
	}

	return token, nil
}
