package model

import (
	"time"

	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
)

// Order - Заказ
type Order struct {
	ID         uint64              `db:"id"`
	Number     string              `db:"number"`
	CustomerID uint64              `db:"customer_id"`
	Status     statusInner.Accrual `db:"status"`
	Accrual    float64             `db:"accrual"` // Начисленные баллы за этот конкретный заказ
	UploadedAt time.Time           `db:"uploaded_at"`
	UpdatedAt  time.Time           `db:"updated_at"`
}

func (o *Order) IsCorrect() bool {
	if o == nil || o.Number == "" || o.CustomerID == 0 || !o.Status.Valid() {
		return false
	}
	return true
}
