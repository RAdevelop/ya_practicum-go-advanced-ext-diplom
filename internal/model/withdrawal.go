package model

import "time"

// Withdrawal - Списание
type Withdrawal struct {
	ID          uint64    `json:"id"`
	CustomerID  uint64    `json:"customer_id"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`          // Сумма списанных баллов для заказа
	ProcessedAt time.Time `json:"processed_at"` // Когда списано
}
