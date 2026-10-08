package model

import statusOuter "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/outer"

// Accrual - Начисление от внешней системы
type Accrual struct {
	Order   string              // Номер заказа
	Status  statusOuter.Accrual // Статус начисления
	Accrual float64             // Количество начисленных баллов
}
