package model

import "time"

// Withdrawal - Списание
type Withdrawal struct {
	ID          uint64    `db:"id"`
	CustomerID  uint64    `db:"customer_id"`
	OrderID     uint64    `db:"order_id"`     // Идентификатор заказа
	Order       string    `db:"order"`        // Номер заказа
	Sum         float64   `db:"sum"`          // Сумма списанных баллов для заказа
	ProcessedAt time.Time `db:"processed_at"` // Когда списано
}
