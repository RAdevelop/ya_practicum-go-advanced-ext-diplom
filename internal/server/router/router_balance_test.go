package router

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetBalance(t *testing.T) {
	t.Skip("TODO implement")

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name:  "StatusOK",
			given: given{},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: `TODO`,
			},
		},
		{
			name:  "StatusUnauthorized",
			given: given{},
			want: want{
				httpStatus:   http.StatusUnauthorized,
				responseBody: ``,
			},
		},
		{
			name:  "StatusInternalServerError",
			given: given{},
			want: want{
				httpStatus:   http.StatusInternalServerError,
				responseBody: ``,
			},
		},
	}

	client := setupServer(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true)
			result, err = req.Get("/api/user/balance")
			assert.NoError(t, err)

			assert.Equalf(t, tt.want.httpStatus, result.StatusCode(), "tt.given: %v", tt.given)
			assert.Equalf(t, "application/json", result.Header().Get("Content-Type"), "tt.given: %v", tt.given)

			// io.Discard выступает в качестве приёмника ненужных данных
			_, err = io.Copy(io.Discard, result.RawResponse.Body)
			assert.NoError(t, err)
			assert.NoError(t, result.RawResponse.Body.Close())
		})
	}
}
func Test_PostBalanceWithdraw(t *testing.T) {
	t.Skip("TODO implement")
}

func Test_GetWithdrawals(t *testing.T) {
	t.Skip("TODO implement")
}
