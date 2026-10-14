package dto

type OrderAccrualQueue struct {
	OrderNumber string `json:"order_number"`
	/*TODO ? CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	NextCheckAt time.Time `json:"next_check_at"`
	Attempts    uint      `json:"attempts"`*/
}
