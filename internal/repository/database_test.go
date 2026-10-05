package repository

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/repository/database"
	configDB "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/repository/database/config"
	"github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/assert"
)

var envOpts = &env.Options{
	Environment: map[string]string{
		"DATABASE_URI": os.Getenv("DB_DSN_TEST"),
	},
}

var errForTransactionRollback = errors.New("return for transaction rollback")

func setUpStorage(t *testing.T) *DBStorage {

	envDB, err := configDB.NewEnvWithOptions(envOpts)
	assert.NoError(t, err)

	ctx := context.Background()
	db, err := database.NewDB(ctx, envDB, nil)
	assert.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})

	return NewDBStorage(db)
}
