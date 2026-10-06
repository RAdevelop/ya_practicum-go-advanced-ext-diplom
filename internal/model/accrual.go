package model

import (
	statusOuter "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/outer"
)

// Accrual - Начисление от внешней системы TODO перенести в dto?! так как такие данные получаем от внешней системы, и обновляем данные в orders!?
type Accrual struct {
	Order   string              `json:"order"`
	Status  statusOuter.Accrual `json:"status"`            // Статус расчёта начисления
	Accrual float64             `json:"accrual,omitempty"` // Рассчитанные баллы к начислению
}
