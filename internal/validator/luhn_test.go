package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidLuhn(t *testing.T) {

	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{
			name:   "valid number",
			number: "12345678903",
			want:   true,
		},
		{
			name:   "not valid number",
			number: "123456",
			want:   false,
		},
		{
			name:   "is string",
			number: "string",
			want:   false,
		},
		{
			name:   "empty string",
			number: "",
			want:   false,
		},
		{
			name:   "mixed string number",
			number: "st12345678903ing",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsValidLuhn(tt.number), "IsValidLuhn(%s)", tt.number)
		})
	}
}
