package database

import (
	"context"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
)

func (s *DBStorage) OrderUpload(ctx context.Context, customerID uint64, number string) error {

	//TODO implement
	return nil
}

func (s *DBStorage) Orders(ctx context.Context, customerID uint64) ([]model.Order, error) {
	//TODO implement
	return nil, nil
}
