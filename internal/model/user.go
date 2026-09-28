package model

import "time"

// User - Пользователь
type User struct {
	ID           uint64    `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
