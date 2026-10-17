package outer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccrual_String(t *testing.T) {
	tests := []struct {
		status Accrual
		want   string
	}{
		{AccrualRegistered, "REGISTERED"},
		{AccrualInvalid, "INVALID"},
		{AccrualProcessing, "PROCESSING"},
		{AccrualProcessed, "PROCESSED"},
		{Accrual{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.String())
		})
	}
}

func TestAccrual_Valid(t *testing.T) {
	tests := []struct {
		name   string
		status Accrual
		want   bool
	}{
		{"REGISTERED", AccrualRegistered, true},
		{"INVALID", AccrualInvalid, true},
		{"PROCESSING", AccrualProcessing, true},
		{"PROCESSED", AccrualProcessed, true},
		{"empty", Accrual{}, false},
		{"unknown", Accrual{value: "UNKNOWN"}, false},
		{"lowercase", Accrual{value: "registered"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.Valid())
		})
	}
}

func TestAccrual_IsFinal(t *testing.T) {
	tests := []struct {
		name   string
		status Accrual
		want   bool
	}{
		{"REGISTERED", AccrualRegistered, false},
		{"PROCESSING", AccrualProcessing, false},
		{"INVALID", AccrualInvalid, true},
		{"PROCESSED", AccrualProcessed, true},
		{"empty", Accrual{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.status
			assert.Equal(t, tt.want, s.IsFinal())
		})
	}
}

func TestAccrual_MarshalJSON(t *testing.T) {
	tests := []struct {
		status Accrual
		want   string
	}{
		{AccrualRegistered, `"REGISTERED"`},
		{AccrualInvalid, `"INVALID"`},
		{AccrualProcessing, `"PROCESSING"`},
		{AccrualProcessed, `"PROCESSED"`},
		{Accrual{}, `""`},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			s := tt.status
			data, err := s.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(data))
		})
	}
}

func TestAccrual_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    Accrual
		wantErr bool
	}{
		{"REGISTERED", `"REGISTERED"`, AccrualRegistered, false},
		{"PROCESSED", `"PROCESSED"`, AccrualProcessed, false},
		{"lowercase", `"processed"`, AccrualProcessed, false},
		{"mixed case", `"ProCessed"`, AccrualProcessed, false},
		{"with spaces", `"  REGISTERED  "`, AccrualRegistered, false},
		{"empty", `""`, Accrual{}, true},
		{"unknown", `"UNKNOWN"`, Accrual{}, true},
		{"number", `123`, Accrual{}, true},
		{"null", `null`, Accrual{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Accrual
			err := s.UnmarshalJSON([]byte(tt.raw))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, s)
		})
	}
}

func TestAccrual_JSONRoundTrip(t *testing.T) {
	statuses := []Accrual{
		AccrualRegistered,
		AccrualInvalid,
		AccrualProcessing,
		AccrualProcessed,
	}

	for _, original := range statuses {
		t.Run(original.String(), func(t *testing.T) {
			data, err := original.MarshalJSON()
			assert.NoError(t, err)

			var parsed Accrual
			err = parsed.UnmarshalJSON(data)
			assert.NoError(t, err)

			assert.Equal(t, original, parsed)
		})
	}
}
