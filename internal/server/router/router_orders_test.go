package router

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// GET /api/user/orders
func Test_GetOrders(t *testing.T) {

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
				responseBody: `TODO`, //TODO
			},
		},
		{
			name:  "StatusNoContent",
			given: given{},
			want: want{
				httpStatus:   http.StatusNoContent,
				responseBody: ``,
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

			result, err = req.Get("/api/user/orders")
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

func Test_PostOrders(t *testing.T) {
	t.Skip("TODO implement")
}
