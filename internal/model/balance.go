package model

// Balance - Баланс
type Balance struct {
	ID         uint64  `json:"id"`
	CustomerID uint64  `json:"customer_id"`
	Current    float64 `json:"current"`   // Текущий баланс
	Withdrawn  float64 `json:"withdrawn"` // Всего списано
}
