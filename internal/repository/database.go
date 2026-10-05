package repository

import (
	"errors"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/repository/database"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBStorage struct {
	DB *database.DB
}

func NewDBStorage(db *database.DB) *DBStorage {
	return &DBStorage{
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
