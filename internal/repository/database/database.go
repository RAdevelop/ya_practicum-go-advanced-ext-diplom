package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type Storage struct {
	DB *DB
}

func NewStorage(db *DB) *Storage {
	return &Storage{
		DB: db,
	}
}

const (
	// pgErrUniqueViolationCode - запись уже существует
	pgErrUniqueViolationCode = "23505"
)

// isPgErrorCode проверяет, что ошибка — PgError с заданным кодом.
func isPgErrorCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}
