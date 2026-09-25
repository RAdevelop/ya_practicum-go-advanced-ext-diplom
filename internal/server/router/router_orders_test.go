package router

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/service"
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	"github.com/stretchr/testify/assert"
)

func Test_GetOrders(t *testing.T) {

	type given struct {
		userID     uint64
		orderList  []model.Order
		orderError error
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "StatusOK",
			given: given{
				userID: 1,
				orderList: []model.Order{
					{
						Number:     "12345",
						UserID:     1,
						Status:     statusInner.AccrualNew,
						UploadedAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
					},
				},
				orderError: nil,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: `[{"number":"12345","user_id":1,"status":"NEW","uploaded_at":"2026-09-25T13:12:16Z"}]`,
				contentType:  "application/json",
			},
		},
		{
			name: "StatusNoContent",
			given: given{
				userID:     1,
				orderList:  nil,
				orderError: nil,
			},
			want: want{
				httpStatus:   http.StatusNoContent,
				responseBody: ``,
				contentType:  "application/json",
			},
		},
		/*{
			TODO name:  "StatusUnauthorized",
			given: given{},
			want: want{
				httpStatus:   http.StatusUnauthorized,
				responseBody: ``,
			},
		},*/
		{
			name: "StatusInternalServerError",
			given: given{
				userID:     1,
				orderList:  nil,
				orderError: errors.New("some error"),
			},
			want: want{
				httpStatus:   http.StatusInternalServerError,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			loyaltyStorage := service.NewMockLoyaltyStorage(t)
			loyaltyStorage.EXPECT().Orders(tt.given.userID).Maybe().Return(tt.given.orderList, tt.given.orderError)

			client := setupServer(t, loyaltyStorage)

			req := client.R().
				SetHeader("Content-Type", "application/json").SetDoNotParseResponse(true)

			result, err = req.Get(uriApiUserOrders)
			assert.NoError(t, err)

			assertResult(t, result, tt.want, tt.given)
		})
	}
}

func Test_PostOrders(t *testing.T) {
	t.Skip("TODO implement")
}
