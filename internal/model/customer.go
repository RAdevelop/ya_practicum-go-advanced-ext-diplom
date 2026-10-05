package model

import (
	"strings"
	"time"
)

// Customer - Пользователь/покупатель
type Customer struct {
	ID           uint64    `json:"id" db:"id"`
	Login        string    `json:"login" db:"login"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

func (c *Customer) IsCorrect() bool {
	if c == nil || strings.TrimSpace(c.Login) == "" || strings.TrimSpace(c.PasswordHash) == "" {
		return false
	}
	return true
}
