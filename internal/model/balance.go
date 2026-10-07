package model

// Balance - Баланс
type Balance struct {
	ID         uint64  `db:"id"`
	CustomerID uint64  `db:"customer_id"`
	Current    float64 `db:"current"`   // Текущий баланс
	Withdrawn  float64 `db:"withdrawn"` // Всего списано
}
