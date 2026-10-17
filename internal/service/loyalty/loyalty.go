package loyalty

import (
	"context"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
)

// LoyaltyStorage - интерфейс для определения работы с хранилищем данных по программе лояльности
//
//go:generate mockery
type LoyaltyStorage interface {
	// OrderUpload - создание нового заказа
	OrderUpload(ctx context.Context, order model.Order) error
	// OrdersByCustomerID - получение списка заказов указанного покупателя
	OrdersByCustomerID(ctx context.Context, customerID uint64) ([]model.Order, error)
	// OrderAccrualQueue - заказы, ожидающие начисления
	OrderAccrualQueue(ctx context.Context, nextCheckAtLessOrEq time.Time, limit uint) ([]model.OrderAccrualQueue, error)
	// OrderAccrualQueueUpdate - обновляем состояние очереди для заказов
	OrderAccrualQueueUpdate(ctx context.Context, orderForUpdate []string, orderForDelete []string) error
	// BalanceByCustomerID - Получение текущего баланса пользователя
	BalanceByCustomerID(ctx context.Context, customerID uint64) (*model.Balance, error)
	// BalanceWithdrawalsByCustomerID - Получение информации о выводе средств
	BalanceWithdrawalsByCustomerID(ctx context.Context, customerID uint64) ([]model.Withdrawal, error)
	// BalanceWithdraw - списание средств с баланса покупателя
	BalanceWithdraw(ctx context.Context, withdrawal model.Withdrawal) error
	// BalanceAccrual - начисление баллов к заказу
	BalanceAccrual(ctx context.Context, accrual model.Accrual) error
	//CustomerCreate - создание нового покупателя
	CustomerCreate(ctx context.Context, customer *model.Customer) (*model.Customer, error)
	// CustomerFindByLogin - поиск покупателя по логину
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
