package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithdrawal_IsCorrect(t *testing.T) {
	tests := []struct {
		name       string
		withdrawal *Withdrawal
		want       bool
	}{
		{
			name:       "nil withdrawal",
			withdrawal: nil,
			want:       false,
		},
		{
			name:       "empty struct",
			withdrawal: &Withdrawal{},
			want:       false,
		},
		{
			name: "empty order number",
			withdrawal: &Withdrawal{
				OrderNumber: "",
				CustomerID:  1,
				Sum:         100,
			},
			want: false,
		},
		{
			name: "order number only spaces",
			withdrawal: &Withdrawal{
				OrderNumber: "   ",
				CustomerID:  1,
				Sum:         100,
			},
			want: false,
		},
		{
			name: "order number tabs and newlines",
			withdrawal: &Withdrawal{
				OrderNumber: "\t\n ",
				CustomerID:  1,
				Sum:         100,
			},
			want: false,
		},
		{
			name: "zero customer id",
			withdrawal: &Withdrawal{
				OrderNumber: "12345",
				CustomerID:  0,
				Sum:         100,
			},
			want: false,
		},
		{
			name: "zero sum",
			withdrawal: &Withdrawal{
				OrderNumber: "12345",
				CustomerID:  1,
				Sum:         0,
			},
			want: false,
		},
		{
			name: "negative sum",
			withdrawal: &Withdrawal{
				OrderNumber: "12345",
				CustomerID:  1,
				Sum:         -10,
			},
			want: false,
		},
		{
			name: "valid withdrawal",
			withdrawal: &Withdrawal{
				OrderNumber: "12345",
				CustomerID:  1,
				Sum:         100,
			},
			want: true,
		},
		{
			name: "valid with leading and trailing spaces in order number",
			withdrawal: &Withdrawal{
				OrderNumber: "  12345  ",
				CustomerID:  1,
				Sum:         100,
			},
			want: true,
		},
		{
			name: "valid with fractional sum",
			withdrawal: &Withdrawal{
				OrderNumber: "12345",
				CustomerID:  1,
				Sum:         0.5,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.withdrawal.IsCorrect())
		})
	}
}
