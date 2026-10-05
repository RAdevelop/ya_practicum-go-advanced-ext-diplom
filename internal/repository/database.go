package repository

import (
	"context"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/repository/database"
)

type DBStorage struct {
	DB *database.DB
}

func NewDBStorage(db *database.DB) *DBStorage {
	return &DBStorage{
		DB: db,
	}
}

func (s *DBStorage) OrderUpload(ctx context.Context, userID uint64, number string) error {

	//TODO implement
	return nil
}

func (s *DBStorage) Orders(ctx context.Context, userID uint64) ([]model.Order, error) {
	//TODO implement
	return nil, nil
}

func (s *DBStorage) Balance(ctx context.Context, userID uint64) (*model.Balance, error) {
	//TODO implement
	return nil, nil
}

func (s *DBStorage) BalanceWithdrawals(ctx context.Context, userID uint64) ([]model.Withdrawal, error) {
	//TODO implement
	return nil, nil
}
func (s *DBStorage) BalanceWithdraw(ctx context.Context, userID uint64, orderNumber string, sum float64) error {
	//TODO implement
	return nil
}
func (s *DBStorage) UserCreate(ctx context.Context, user *model.User) (*model.User, error) {
	//TODO implement
	return nil, nil
}
func (s *DBStorage) UserByLogin(ctx context.Context, login string) (*model.User, error) {
	//TODO implement
	return nil, nil
}
