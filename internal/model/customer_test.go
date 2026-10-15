package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomer_IsCorrect(t *testing.T) {
	tests := []struct {
		name     string
		customer *Customer
		want     bool
	}{
		{
			name:     "nil customer",
			customer: nil,
			want:     false,
		},
		{
			name:     "empty struct",
			customer: &Customer{},
			want:     false,
		},
		{
			name: "empty login",
			customer: &Customer{
				Login:        "",
				PasswordHash: "hash",
			},
			want: false,
		},
		{
			name: "empty password hash",
			customer: &Customer{
				Login:        "alice",
				PasswordHash: "",
			},
			want: false,
		},
		{
			name: "login only spaces",
			customer: &Customer{
				Login:        "   ",
				PasswordHash: "hash",
			},
			want: false,
		},
		{
			name: "password hash only spaces",
			customer: &Customer{
				Login:        "alice",
				PasswordHash: "   ",
			},
			want: false,
		},
		{
			name: "login with tabs and newlines",
			customer: &Customer{
				Login:        "\t\n ",
				PasswordHash: "hash",
			},
			want: false,
		},
		{
			name: "valid customer",
			customer: &Customer{
				ID:           1,
				Login:        "alice",
				PasswordHash: "hash",
			},
			want: true,
		},
		{
			name: "valid with leading/trailing spaces in login",
			customer: &Customer{
				Login:        "  alice  ",
				PasswordHash: "hash",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.customer.IsCorrect())
		})
	}
}
