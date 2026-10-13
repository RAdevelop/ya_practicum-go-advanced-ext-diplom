package model

import "time"

// Withdrawal - Списание
type Withdrawal struct {
	OrderNumber string    `db:"order_number"` // Номер заказа
	CustomerID  uint64    `db:"customer_id"`
	Sum         float64   `db:"sum"`          // Сумма списанных баллов для заказа
	ProcessedAt time.Time `db:"processed_at"` // Когда списано
}

func (w *Withdrawal) IsCorrect() bool {
	if w == nil || w.OrderNumber == "" || w.CustomerID == 0 || w.Sum <= 0 {
		return false
	}
	return true
}
