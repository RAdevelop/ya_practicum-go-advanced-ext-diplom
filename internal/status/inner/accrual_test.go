package inner

import (
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	statusOuter "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/outer"
	"github.com/stretchr/testify/assert"
)

func TestAccrual_String(t *testing.T) {
	tests := []struct {
		status Accrual
		want   string
	}{
		{AccrualNew, "NEW"},
		{AccrualProcessing, "PROCESSING"},
		{AccrualInvalid, "INVALID"},
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
		{"NEW", AccrualNew, true},
		{"PROCESSING", AccrualProcessing, true},
		{"INVALID", AccrualInvalid, true},
		{"PROCESSED", AccrualProcessed, true},
		{"empty", Accrual{}, false},
		{"unknown", Accrual{value: "UNKNOWN"}, false},
		{"lowercase", Accrual{value: "new"}, false},
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
		{"NEW", AccrualNew, false},
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
		{AccrualNew, `"NEW"`},
		{AccrualProcessing, `"PROCESSING"`},
		{AccrualInvalid, `"INVALID"`},
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
		wantErr error
	}{
		{"NEW", `"NEW"`, AccrualNew, nil},
		{"PROCESSED", `"PROCESSED"`, AccrualProcessed, nil},
		{"lowercase", `"processed"`, AccrualProcessed, nil},
		{"with spaces", `"  NEW  "`, AccrualNew, nil},
		{"mixed case with spaces", `" ProCessed "`, AccrualProcessed, nil},
		{"empty", `""`, Accrual{}, perror.ErrAccrualStatusInvalid},
		{"unknown", `"UNKNOWN"`, Accrual{}, perror.ErrAccrualStatusInvalid},
		{"not a string", `123`, Accrual{}, nil}, // ошибка от json.Unmarshal
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Accrual
			err := s.UnmarshalJSON([]byte(tt.raw))

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			if tt.raw == "123" {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, s)
		})
	}
}

func TestAccrual_Scan(t *testing.T) {
	tests := []struct {
		name    string
		src     any
		want    Accrual
		wantErr error
	}{
		{"string NEW", "NEW", AccrualNew, nil},
		{"string processed", "processed", AccrualProcessed, nil},
		{"string with spaces", "  PROCESSING  ", AccrualProcessing, nil},
		{"bytes NEW", []byte("NEW"), AccrualNew, nil},
		{"bytes invalid", []byte("INVALID"), AccrualInvalid, nil},
		{"nil", nil, Accrual{}, perror.ErrAccrualStatusInvalid},
		{"int", 123, Accrual{}, perror.ErrAccrualStatusInvalid},
		{"unknown", "UNKNOWN", Accrual{}, perror.ErrAccrualStatusInvalid},
		{"empty string", "", Accrual{}, perror.ErrAccrualStatusInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Accrual
			err := s.Scan(tt.src)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, s)
		})
	}
}

func TestAccrual_ConvertFromOuter(t *testing.T) {
	tests := []struct {
		name  string
		outer statusOuter.Accrual
		want  Accrual
	}{
		{"REGISTERED -> NEW", statusOuter.AccrualRegistered, AccrualNew},
		{"PROCESSING -> PROCESSING", statusOuter.AccrualProcessing, AccrualProcessing},
		{"INVALID -> INVALID", statusOuter.AccrualInvalid, AccrualInvalid},
		{"PROCESSED -> PROCESSED", statusOuter.AccrualProcessed, AccrualProcessed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Accrual
			s.ConvertFromOuter(tt.outer)
			assert.Equal(t, tt.want, s)
		})
	}
}
