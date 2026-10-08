package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) CustomerCreate(ctx context.Context, customer *model.Customer) (*model.Customer, error) {

	if !customer.IsCorrect() {
		return nil, perror.ErrCustomerInvalidModel
	}

	const sql = `
		INSERT INTO customers (login, password_hash)
		VALUES ($1, $2)
		RETURNING id, login, password_hash, created_at
		`

	row, err := s.DB.Executor(ctx).Query(ctx, sql, customer.Login, customer.PasswordHash)

	customer, err = pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.Customer])

	if isPgErrorCode(err, pgErrUniqueViolationCode) {
		return nil, errors.Join(perror.ErrCustomerAlreadyExists, err)
	}

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *Storage) CustomerFindByLogin(ctx context.Context, login string) (*model.Customer, error) {

	if strings.TrimSpace(login) == "" {
		return nil, perror.ErrCustomerNotFound
	}

	const sql = `SELECT id, login, password_hash, created_at FROM customers WHERE login = $1`

	row, err := s.DB.Executor(ctx).Query(ctx, sql, login)

	if err != nil {
		return nil, fmt.Errorf("%w, %w", perror.ErrCustomerNotFound, err)
	}

	customer, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.Customer])
	if err != nil {
		return nil, fmt.Errorf("%w, %w", perror.ErrCustomerNotFound, err)
	}
	return customer, nil
}
