package database

import (
	"context"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
)

/*
TODO учитывать:
 - помнить о транзакциях и/или блокировка записей перед начислением/списанием баллов
*/

func (s *Storage) BalanceByCustomerID(ctx context.Context, customerID uint64) (*model.Balance, error) {
	//TODO implement
	return nil, nil
}

func (s *Storage) BalanceWithdrawalsByCustomerID(ctx context.Context, customerID uint64) ([]model.Withdrawal, error) {
	//TODO implement
	return nil, nil
}
func (s *Storage) BalanceWithdraw(ctx context.Context, customerID uint64, orderNumber string, sum float64) error {
	//TODO implement
	return nil
}
