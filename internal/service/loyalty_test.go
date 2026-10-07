package service

import (
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

const testJWTSecret = "lPD7WBZ/MCBKK0aEqgzSqfIQSqAGB7VhIfjsZwuXLJE="

var testCustomerDTO = &dto.Customer{
	ID:    1,
	Login: "TestLogin",
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	return string(hash)
}
