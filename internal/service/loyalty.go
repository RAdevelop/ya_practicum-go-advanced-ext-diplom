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

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/retryer"
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
func (lm *LoyaltyManager) OrderUpload(ctx context.Context, userID uint64, number string) error {
	var err error
	_, err = retryer.RetryLinear(ctx, func(ctx context.Context) (struct{}, error) {
		err = lm.storage.OrderUpload(ctx, userID, number)
		return struct{}{}, err
	}, retryLinearStepSeconds, new(retryLinearAttempts))

	return err
}

// Orders - Получение списка загруженных номеров заказов
func (lm *LoyaltyManager) Orders(ctx context.Context, userID uint64) ([]model.Order, error) {
	return retryer.RetryLinear(ctx, func(ctx context.Context) ([]model.Order, error) {
		return lm.storage.Orders(ctx, userID)
	}, retryLinearStepSeconds, new(retryLinearAttempts))
}

// Balance - Получение текущего баланса пользователя
func (lm *LoyaltyManager) Balance(ctx context.Context, userID uint64) (*model.Balance, error) {
	return retryer.RetryLinear(ctx, func(ctx context.Context) (*model.Balance, error) {
		return lm.storage.Balance(ctx, userID)
	}, retryLinearStepSeconds, new(retryLinearAttempts))
}

// BalanceWithdrawals - Получение информации о выводе средств
func (lm *LoyaltyManager) BalanceWithdrawals(ctx context.Context, userID uint64) ([]model.Withdrawal, error) {
	return retryer.RetryLinear(ctx, func(ctx context.Context) ([]model.Withdrawal, error) {
		return lm.storage.BalanceWithdrawals(ctx, userID)
	}, retryLinearStepSeconds, new(retryLinearAttempts))
}

// BalanceWithdraw - списание средств
func (lm *LoyaltyManager) BalanceWithdraw(ctx context.Context, userID uint64, orderNumber string, sum float64) error {

	var err error
	_, err = retryer.RetryLinear(ctx, func(ctx context.Context) (struct{}, error) {
		err = lm.storage.BalanceWithdraw(ctx, userID, orderNumber, sum)
		return struct{}{}, err
	}, retryLinearStepSeconds, new(retryLinearAttempts))

	return err
}
