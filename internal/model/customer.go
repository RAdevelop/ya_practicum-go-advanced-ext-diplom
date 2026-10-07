package model

import (
	"strings"
	"time"
)

// Customer - Пользователь/покупатель
type Customer struct {
	ID           uint64    `db:"id"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

func (c *Customer) IsCorrect() bool {
	if c == nil || strings.TrimSpace(c.Login) == "" || strings.TrimSpace(c.PasswordHash) == "" {
		return false
	}
	return true
}
