// Package outer - для описания статусов внешней системы
package outer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Accrual - статус расчёта начисления от внешнего сервиса
type Accrual struct {
	value string
}

var (
	AccrualRegistered = Accrual{"REGISTERED"} // Заказ зарегистрирован, но начисление не рассчитано
	AccrualInvalid    = Accrual{"INVALID"}    // Заказ не принят к расчёту, и вознаграждение не будет начислено
	AccrualProcessing = Accrual{"PROCESSING"} // Расчёт начисления в процессе
	AccrualProcessed  = Accrual{"PROCESSED"}  // Расчёт начисления окончен
)

func (s Accrual) String() string {
	return s.value
}

func (s Accrual) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.value)
}

func (s *Accrual) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	normalized := Accrual{strings.ToUpper(strings.TrimSpace(raw))}
	if !normalized.valid() {
		return fmt.Errorf("invalid status: %s", raw)
	}
	return nil
}

func (s Accrual) valid() bool {
	switch s.value {
	case AccrualRegistered.value, AccrualInvalid.value, AccrualProcessing.value, AccrualProcessed.value:
		return true
	default:
		return false
	}
}

func (s Accrual) IsFinal() bool {
	return s == AccrualInvalid || s == AccrualProcessed
}
