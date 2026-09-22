package model

import (
	"time"

	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
)

// Order - Заказ
type Order struct {
	Number     string              `json:"number"`
	UserID     uint64              `json:"user_id"`
	Status     statusInner.Accrual `json:"status"`
	Accrual    float64             `json:"accrual,omitempty"` // Начисленные баллы за этот конкретный заказ
	UploadedAt time.Time           `json:"uploaded_at"`
}
