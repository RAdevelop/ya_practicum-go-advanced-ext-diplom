// Package inner - для описания статусов внутренней системы
package inner

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
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
func (s *Accrual) String() string {
	return s.value
}

// MarshalJSON сериализует статус как JSON-строку.
func (s *Accrual) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.value)
}

// UnmarshalJSON десериализует статус из JSON-строки с валидацией.
func (s *Accrual) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	normalized := normalize(raw)
	if !normalized.Valid() {
		return fmt.Errorf("%w: %q", perror.ErrAccrualStatusInvalid, raw)
	}

	*s = normalized
	return nil
}

// IsFinal - сообщает, является ли статус окончательным.
func (s *Accrual) IsFinal() bool {
	return *s == AccrualInvalid || *s == AccrualProcessed
}

// Valid - проверяет, что статус — одно из известных значений.
func (s *Accrual) Valid() bool {
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

// Scan - реализует sql.Scanner — читает строку из БД.
func (s *Accrual) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("%w: %T", perror.ErrAccrualStatusInvalid, src)
	}

	var raw string
	switch v := src.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("%w: %T", perror.ErrAccrualStatusInvalid, src)
	}

	normalized := normalize(raw)
	if !normalized.Valid() {
		return fmt.Errorf("%w: %q", perror.ErrAccrualStatusInvalid, raw)
	}

	*s = normalized
	return nil
}

func normalize(raw string) Accrual {
	return Accrual{value: strings.ToUpper(strings.TrimSpace(raw))}
}
