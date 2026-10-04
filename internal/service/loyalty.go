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
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/retryer"
	"golang.org/x/crypto/bcrypt"
)

// LoyaltyStorage - интерфейс для определения работы с хранилищем данных по программе лояльности
//
//go:generate mockery
type LoyaltyStorage interface {
	OrderUpload(ctx context.Context, userID uint64, number string) error
	Orders(ctx context.Context, userID uint64) ([]model.Order, error)
	Balance(ctx context.Context, userID uint64) (*model.Balance, error)
	BalanceWithdrawals(ctx context.Context, userID uint64) ([]model.Withdrawal, error)
	BalanceWithdraw(ctx context.Context, userID uint64, orderNumber string, sum float64) error
	UserCreate(ctx context.Context, user *model.User) (*model.User, error)
	UserByLogin(ctx context.Context, login string) (*model.User, error)
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

const retryLinearStepSeconds uint = 2
const retryLinearAttempts uint = 3

// OrderUpload - Загрузка заказа
func (lm *LoyaltyManager) OrderUpload(ctx context.Context, userDTO *dto.User, number string) error {
	var err error
	_, err = retryer.RetryLinear(ctx, func(ctx context.Context) (struct{}, error) {
		err = lm.storage.OrderUpload(ctx, userDTO.ID, number)
		return struct{}{}, err
	}, retryLinearStepSeconds, new(retryLinearAttempts))

	return err
}

// Orders - Получение списка загруженных номеров заказов
func (lm *LoyaltyManager) Orders(ctx context.Context, userDTO *dto.User) ([]model.Order, error) {
	return retryer.RetryLinear(ctx, func(ctx context.Context) ([]model.Order, error) {
		return lm.storage.Orders(ctx, userDTO.ID)
	}, retryLinearStepSeconds, new(retryLinearAttempts))
}

// Balance - Получение текущего баланса пользователя
func (lm *LoyaltyManager) Balance(ctx context.Context, userDTO *dto.User) (*model.Balance, error) {
	return retryer.RetryLinear(ctx, func(ctx context.Context) (*model.Balance, error) {
		return lm.storage.Balance(ctx, userDTO.ID)
	}, retryLinearStepSeconds, new(retryLinearAttempts))
}

// BalanceWithdrawals - Получение информации о выводе средств
func (lm *LoyaltyManager) BalanceWithdrawals(ctx context.Context, userDTO *dto.User) ([]model.Withdrawal, error) {
	return retryer.RetryLinear(ctx, func(ctx context.Context) ([]model.Withdrawal, error) {
		return lm.storage.BalanceWithdrawals(ctx, userDTO.ID)
	}, retryLinearStepSeconds, new(retryLinearAttempts))
}

// BalanceWithdraw - списание средств
func (lm *LoyaltyManager) BalanceWithdraw(ctx context.Context, userDTO *dto.User, orderNumber string, sum float64) error {

	var err error
	_, err = retryer.RetryLinear(ctx, func(ctx context.Context) (struct{}, error) {
		err = lm.storage.BalanceWithdraw(ctx, userDTO.ID, orderNumber, sum)
		return struct{}{}, err
	}, retryLinearStepSeconds, new(retryLinearAttempts))

	return err
}

func (lm *LoyaltyManager) UserLogin(ctx context.Context, userCredentials dto.UserCredentials, jwtSecret string) (string, error) {

	userModel, err := retryer.RetryLinear(ctx, func(ctx context.Context) (*model.User, error) {
		return lm.storage.UserByLogin(ctx, userCredentials.Login)
	}, retryLinearStepSeconds, new(retryLinearAttempts))

	if err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrUserNotFound, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userModel.PasswordHash), []byte(userCredentials.Password)); err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrInvalidUserCredentials, err)
	}

	userDTO := dto.User{
		ID:    userModel.ID,
		Login: userModel.Login,
	}
	token, err := jwtoken.Generate(userDTO, []byte(jwtSecret))

	if err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrInvalidToken, err)
	}

	return token, nil
}

func (lm *LoyaltyManager) UserRegister(ctx context.Context, userCredentials dto.UserCredentials, jwtSecret string) (string, error) {

	// 1. Хешируем пароль.
	hash, err := bcrypt.GenerateFromPassword([]byte(userCredentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrHashGenerateFromPassword, err)
	}

	// 2. Создаём пользователя.
	userModel := &model.User{
		Login:        userCredentials.Login,
		PasswordHash: string(hash),
	}

	userModel, err = retryer.RetryLinear(ctx, func(ctx context.Context) (*model.User, error) {
		return lm.storage.UserCreate(ctx, userModel)
	}, retryLinearStepSeconds, new(retryLinearAttempts))

	if err != nil {
		return "", err
	}

	/*
		Можно считать, что идет дублирование данных.
		Зато слой сервиса генерации токена ничего не знает именно о модели пользователя.
		И в этом случае необходимый рефакторинг слоев можно будет делать независимо.
	*/
	userDTO := dto.User{
		Login: userModel.Login,
		ID:    userModel.ID,
	}

	// 3. выдаём JWT.
	token, err := jwtoken.Generate(userDTO, []byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrInvalidToken, err)
	}

	return token, nil
}
