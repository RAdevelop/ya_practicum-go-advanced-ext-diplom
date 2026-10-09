package model

import "time"

// Withdrawal - Списание
type Withdrawal struct {
	OrderID     uint64    `db:"order_id"` // Идентификатор заказа
	CustomerID  uint64    `db:"customer_id"`
	Order       string    `db:"order"`        // Номер заказа
	Sum         float64   `db:"sum"`          // Сумма списанных баллов для заказа
	ProcessedAt time.Time `db:"processed_at"` // Когда списано
}

func (w *Withdrawal) IsCorrect() bool {
	if w == nil || w.Order == "" || w.CustomerID == 0 || w.Sum <= 0 {
		return false
	}
	return true
}
