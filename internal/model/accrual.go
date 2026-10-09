package model

import (
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
)

// Accrual - Начисление от внешней системы
type Accrual struct {
	Order   string              // Номер заказа
	Status  statusInner.Accrual // Статус начисления
	Accrual float64             // Количество начисленных баллов
}
