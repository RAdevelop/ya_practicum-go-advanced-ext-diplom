package router

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/service"
	"github.com/stretchr/testify/assert"
)

func Test_GetBalance(t *testing.T) {

	type given struct {
		userID     uint64
		balance    *model.Balance
		balanceErr error
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
				balance: &model.Balance{
					ID:        1,
					UserID:    1,
					Current:   19.05,
					Withdrawn: 2019.05,
				},
				balanceErr: nil,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: `{"id":1,"user_id":1,"current":19.05,"withdrawn":2019.05}`,
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
				balance:    nil,
				balanceErr: errors.New("some error"),
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
			loyaltyStorage.EXPECT().Balance(tt.given.userID).Return(tt.given.balance, tt.given.balanceErr)

			client := setupServer(t, loyaltyStorage)

			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true)
			result, err = req.Get(uriApiUserBalance)
			assert.NoError(t, err)

			assertResult(t, result, tt.want, tt.given)
		})
	}
}
func Test_PostBalanceWithdraw(t *testing.T) {
	t.Skip("TODO implement")
}

func Test_GetWithdrawals(t *testing.T) {

	type given struct {
		userID         uint64
		withdrawals    []model.Withdrawal
		withdrawalsErr error
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
				withdrawals: []model.Withdrawal{
					{
						ID:          123,
						UserID:      1,
						Order:       "12345",
						Sum:         123.45,
						ProcessedAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
					},
				},
				withdrawalsErr: nil,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: `[{"id":123,"user_id":1,"order":"12345","sum":123.45,"processed_at":"2026-09-25T13:12:16Z"}]`,
				contentType:  "application/json",
			},
		},
		{
			name: "StatusNoContent",
			given: given{
				userID:         1,
				withdrawals:    nil,
				withdrawalsErr: nil,
			},
			want: want{
				httpStatus:   http.StatusNoContent,
				responseBody: ``,
				contentType:  "application/json",
			},
		},
		/*{
			TODO name: "StatusUnauthorized",
		},*/
		{
			name: "StatusInternalServerError",
			given: given{
				userID:         1,
				withdrawals:    nil,
				withdrawalsErr: errors.New("some error"),
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
			loyaltyStorage.EXPECT().BalanceWithdrawals(tt.given.userID).Return(tt.given.withdrawals, tt.given.withdrawalsErr)
			client := setupServer(t, loyaltyStorage)
			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true)

			result, err = req.Get(uriApiUserWithdrawals)
			assert.NoErrorf(t, err, "given: %+v", tt.given)

			assertResult(t, result, tt.want, result)
		})
	}
}
