package model

// Balance - Баланс
type Balance struct {
	ID        uint64  `json:"id"`
	UserID    uint64  `json:"user_id"`
	Current   float64 `json:"current"`   // Текущий баланс
	Withdrawn float64 `json:"withdrawn"` // Всего списано
}
