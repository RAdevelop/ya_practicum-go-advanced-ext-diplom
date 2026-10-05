package repository

import (
	"context"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
)

func (s *DBStorage) Balance(ctx context.Context, customerID uint64) (*model.Balance, error) {
	//TODO implement
	return nil, nil
}

func (s *DBStorage) BalanceWithdrawals(ctx context.Context, customerID uint64) ([]model.Withdrawal, error) {
	//TODO implement
	return nil, nil
}
func (s *DBStorage) BalanceWithdraw(ctx context.Context, customerID uint64, orderNumber string, sum float64) error {
	//TODO implement
	return nil
}
