package loyalty

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
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
)

// LoyaltyStorage - интерфейс для определения работы с хранилищем данных по программе лояльности
//
//go:generate mockery
type LoyaltyStorage interface {
	OrderUpload(ctx context.Context, order model.Order) error
	OrdersByCustomerID(ctx context.Context, customerID uint64) ([]model.Order, error)
	OrdersAwaitingAccrual([]statusInner.Accrual) ([]model.Order, error)
	BalanceByCustomerID(ctx context.Context, customerID uint64) (*model.Balance, error)
	BalanceWithdrawalsByCustomerID(ctx context.Context, customerID uint64) ([]model.Withdrawal, error)
	BalanceWithdraw(ctx context.Context, withdrawal model.Withdrawal) error
	BalanceAccrual(ctx context.Context, accrual model.Accrual) error
	CustomerCreate(ctx context.Context, customer *model.Customer) (*model.Customer, error)
	CustomerFindByLogin(ctx context.Context, login string) (*model.Customer, error)
}

// Manager - сервис для работы с программой лояльности
type Manager struct {
	storage LoyaltyStorage
}

func NewManager(storage LoyaltyStorage) *Manager {
	return &Manager{
		storage: storage,
	}
}
