package model

import (
	"testing"

	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	"github.com/stretchr/testify/assert"
)

func TestOrder_IsCorrect(t *testing.T) {
	tests := []struct {
		name  string
		order *Order
		want  bool
	}{
		{
			name:  "nil order",
			order: nil,
			want:  false,
		},
		{
			name:  "empty struct",
			order: &Order{},
			want:  false,
		},
		{
			name: "empty number",
			order: &Order{
				Number:     "",
				CustomerID: 1,
				Status:     statusInner.AccrualNew,
			},
			want: false,
		},
		{
			name: "number only spaces",
			order: &Order{
				Number:     "   ",
				CustomerID: 1,
				Status:     statusInner.AccrualNew,
			},
			want: false,
		},
		{
			name: "number tabs and newlines",
			order: &Order{
				Number:     "\t\n ",
				CustomerID: 1,
				Status:     statusInner.AccrualNew,
			},
			want: false,
		},
		{
			name: "number with leading and trailing spaces",
			order: &Order{
				Number:     "  12345  ",
				CustomerID: 1,
				Status:     statusInner.AccrualNew,
			},
			want: true,
		},
		{
			name: "zero customer id",
			order: &Order{
				Number:     "12345",
				CustomerID: 0,
				Status:     statusInner.AccrualNew,
			},
			want: false,
		},
		{
			name: "invalid status",
			order: &Order{
				Number:     "12345",
				CustomerID: 1,
				Status:     statusInner.Accrual{},
			},
			want: false,
		},
		{
			name: "valid with NEW",
			order: &Order{
				Number:     "12345",
				CustomerID: 1,
				Status:     statusInner.AccrualNew,
			},
			want: true,
		},
		{
			name: "valid with PROCESSING",
			order: &Order{
				Number:     "12345",
				CustomerID: 1,
				Status:     statusInner.AccrualProcessing,
			},
			want: true,
		},
		{
			name: "valid with PROCESSED",
			order: &Order{
				Number:     "12345",
				CustomerID: 1,
				Status:     statusInner.AccrualProcessed,
			},
			want: true,
		},
		{
			name: "valid with INVALID",
			order: &Order{
				Number:     "12345",
				CustomerID: 1,
				Status:     statusInner.AccrualInvalid,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.order.IsCorrect())
		})
	}
}
