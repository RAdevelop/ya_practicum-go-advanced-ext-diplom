package model

import (
	"time"

	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
)

// Order - Заказ
type Order struct {
	ID         uint64              `json:"id" db:"id"`
	Number     string              `json:"number" db:"number"`
	CustomerID uint64              `json:"customer_id" db:"customer_id"`
	Status     statusInner.Accrual `json:"status" db:"status"`
	Accrual    float64             `json:"accrual,omitempty" db:"accrual"` // Начисленные баллы за этот конкретный заказ
	UploadedAt time.Time           `json:"uploaded_at" db:"uploaded_at"`
	UpdatedAt  time.Time           `json:"updated_at" db:"updated_at"`
}

func (o *Order) IsCorrect() bool {
	if o == nil || o.Number == "" || o.CustomerID == 0 || !o.Status.Valid() {
		return false
	}
	return true
}
