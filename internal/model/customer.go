package model

import "time"

// Customer - Пользователь/покупатель
type Customer struct {
	ID           uint64    `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
