package dto

type OrderAccrualQueue struct {
	// Номер заказа
	OrderNumber string
	// Удалить (false) или изменить (true) в очереди
	NeedUpdate bool
}
