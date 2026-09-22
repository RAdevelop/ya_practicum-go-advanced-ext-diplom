// Package inner - для описания статусов внутренней системы
package inner

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Accrual — статус начисления баллов для заказа.
type Accrual struct {
	value string
}

var (
	AccrualNew        = Accrual{"NEW"}        // Заказ загружен в систему, но не попал в обработку
	AccrualProcessing = Accrual{"PROCESSING"} // Вознаграждение за заказ рассчитывается
	AccrualInvalid    = Accrual{"INVALID"}    // Система расчёта вознаграждений отказала в расчёте
	AccrualProcessed  = Accrual{"PROCESSED"}  // Данные по заказу проверены и информация о расчёте успешно получена
)

// String возвращает строковое представление статуса.
func (s Accrual) String() string {
	return s.value
}

// MarshalJSON сериализует статус как JSON-строку.
func (s Accrual) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.value)
}

// UnmarshalJSON десериализует статус из JSON-строки с валидацией.
func (s *Accrual) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	normalized := Accrual{value: strings.ToUpper(strings.TrimSpace(raw))}
	if !normalized.valid() {
		return fmt.Errorf("invalid status: %q", raw)
	}

	*s = normalized
	return nil
}

// IsFinal - сообщает, является ли статус окончательным.
func (s Accrual) IsFinal() bool {
	return s == AccrualInvalid || s == AccrualProcessed
}

// valid - проверяет, что статус — одно из известных значений.
func (s Accrual) valid() bool {
	switch s.value {
	case AccrualNew.value,
		AccrualProcessing.value,
		AccrualInvalid.value,
		AccrualProcessed.value:
		return true
	default:
		return false
	}
}
