package dto

import "time"

// Withdrawal - Списание
type Withdrawal struct {
	Order       string    `json:"order"`        // Номер заказа
	Sum         float64   `json:"sum"`          // Сумма списанных баллов для заказа
	ProcessedAt time.Time `json:"processed_at"` // Когда списано
}
