package dto

import (
	"time"

	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
)

// Order - Заказ
type Order struct {
	Number     string              `json:"number"`
	Status     statusInner.Accrual `json:"status"`
	Accrual    float64             `json:"accrual,omitempty"`
	UploadedAt time.Time           `json:"uploaded_at"`
}
