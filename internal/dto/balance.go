package dto

// Balance - Баланс
type Balance struct {
	Current   float64 `json:"current"`   // Текущий баланс
	Withdrawn float64 `json:"withdrawn"` // Всего списано
}
