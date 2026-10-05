package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *DBStorage) CustomerCreate(ctx context.Context, customer *model.Customer) (*model.Customer, error) {

	if !customer.IsCorrect() {
		return nil, perror.ErrInvalidCustomerModel
	}

	customer.CreatedAt = time.Now()

	sql := `
		INSERT INTO customers (login, password_hash, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, login, password_hash, created_at
		`

	row, err := s.DB.Executor(ctx).Query(ctx, sql, customer.Login, customer.PasswordHash, customer.CreatedAt)

	//TODO del
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		log.Printf("pg error1: code=%s, message=%s, detail=%s, hint=%s",
			pgErr.Code, pgErr.Message, pgErr.Detail, pgErr.Hint)
	}

	customer, err = pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[model.Customer])

	if isPgErrorCode(err, pgErrUniqueViolationCode) {
		return nil, errors.Join(perror.ErrCustomerAlreadyExists, err)
	}

	//TODO del
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		log.Printf("pg error2: code=%s, message=%s, detail=%s, hint=%s",
			pgErr.Code, pgErr.Message, pgErr.Detail, pgErr.Hint)
	}

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *DBStorage) CustomerFindByLogin(ctx context.Context, login string) (*model.Customer, error) {

	if strings.TrimSpace(login) == "" {
		return nil, perror.ErrCustomerNotFound
	}

	sql := `SELECT id, login, password_hash, created_at FROM customers WHERE login = $1`

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
