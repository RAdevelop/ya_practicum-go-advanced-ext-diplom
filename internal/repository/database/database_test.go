package database

import (
	"context"
	"errors"
	"os"
	"testing"

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

func setUpStorage(t *testing.T) *Storage {

	envDB, err := configDB.NewEnvWithOptions(envOpts)
	assert.NoError(t, err)

	ctx := context.Background()
	db, err := NewDB(ctx, envDB, nil)
	assert.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})

	return NewStorage(db)
}
