package router

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
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
			result, err = req.Get(uriUserBalance)
			assert.NoError(t, err)

			assertResult(t, result, tt.want, tt.given)
		})
	}
}

func Test_PostBalanceWithdraw(t *testing.T) {
	type given struct {
		userID      uint64
		orderNumber string
		sum         float64
		balanceErr  error
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "StatusOK",
			given: given{
				userID:      1,
				orderNumber: "4532015112830366",
				sum:         19.05,
				balanceErr:  nil,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		/*{
			TODO name: "StatusUnauthorized",
		},*/
		{
			name: "StatusNotFound",
			given: given{
				userID:      1,
				orderNumber: "4532015112830366",
				sum:         19.05,
				balanceErr:  perror.ErrOrderNotFound,
			},
			want: want{
				httpStatus:   http.StatusNotFound,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusPaymentRequired",
			given: given{
				userID:      1,
				orderNumber: "4532015112830366",
				sum:         19.05,
				balanceErr:  perror.ErrBalanceInsufficient,
			},
			want: want{
				httpStatus:   http.StatusPaymentRequired,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusUnprocessableEntity orderNumber",
			given: given{
				userID:      1,
				orderNumber: "234",
				sum:         19.05,
				balanceErr:  nil,
			},
			want: want{
				httpStatus:   http.StatusUnprocessableEntity,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusUnprocessableEntity sum zero",
			given: given{
				userID:      1,
				orderNumber: "4532015112830366",
				sum:         0,
				balanceErr:  nil,
			},
			want: want{
				httpStatus:   http.StatusUnprocessableEntity,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusUnprocessableEntity sum negative",
			given: given{
				userID:      1,
				orderNumber: "4532015112830366",
				sum:         -100.25,
				balanceErr:  nil,
			},
			want: want{
				httpStatus:   http.StatusUnprocessableEntity,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusInternalServerError",
			given: given{
				userID:      1,
				orderNumber: "4532015112830366",
				sum:         100.25,
				balanceErr:  errors.New("some error"),
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
			loyaltyStorage.EXPECT().BalanceWithdraw(tt.given.userID, tt.given.orderNumber, tt.given.sum).Maybe().Return(tt.given.balanceErr)

			client := setupServer(t, loyaltyStorage)

			balanceWithdraw := dto.BalanceWithdraw{
				OrderNumber: tt.given.orderNumber,
				Sum:         tt.given.sum,
			}

			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true).SetBody(balanceWithdraw)

			result, err = req.Post(uriUserBalanceWithdraw)
			assert.NoError(t, err)
			assertResult(t, result, tt.want, tt.given)
		})
	}
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

			result, err = req.Get(uriUserWithdrawals)
			assert.NoErrorf(t, err, "given: %+v", tt.given)

			assertResult(t, result, tt.want, result)
		})
	}
}
