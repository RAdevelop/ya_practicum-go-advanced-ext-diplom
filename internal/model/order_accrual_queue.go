package model

import "time"

type OrderAccrualQueue struct {
	OrderNumber string    `db:"order_number"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	NextCheckAt time.Time `db:"next_check_at"`
	Attempts    uint      `db:"attempts"`
}
