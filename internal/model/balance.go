package model

import "time"

// Balance - Баланс
type Balance struct {
	CustomerID uint64    `db:"customer_id"`
	Current    float64   `db:"current"`   // Текущий баланс
	Withdrawn  float64   `db:"withdrawn"` // Всего списано
	UpdatedAt  time.Time `db:"updated_at"`
}
