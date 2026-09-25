package service

import "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"

// LoyaltyStorage - интерфейс для определения работы с хранилищем данных по программе лояльности
//
//go:generate mockery
type LoyaltyStorage interface {
	OrderUpload(userID uint64, number string) error
	Orders(userID uint64) ([]model.Order, error)
	Balance(userID uint64) (*model.Balance, error)
	BalanceWithdrawals(userID uint64) ([]model.Withdrawal, error)
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
func (lm *LoyaltyManager) OrderUpload(userID uint64, number string) error {
	return lm.storage.OrderUpload(userID, number)
}

// Orders - Получение списка загруженных номеров заказов
func (lm *LoyaltyManager) Orders(userID uint64) ([]model.Order, error) {
	return lm.storage.Orders(userID)
}

// Balance - Получение текущего баланса пользователя
func (lm *LoyaltyManager) Balance(userID uint64) (*model.Balance, error) {
	return lm.storage.Balance(userID)
}

// BalanceWithdrawals - Получение информации о выводе средств
func (lm *LoyaltyManager) BalanceWithdrawals(userID uint64) ([]model.Withdrawal, error) {
	return lm.storage.BalanceWithdrawals(userID)
}
