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
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetBalance(t *testing.T) {

	type given struct {
		inputUserID        uint64
		outBalance         *model.Balance
		outBalanceErr      error
		makeLoyaltyStorage func(t *testing.T, inputUserID uint64, outBalance *model.Balance, outBalanceErr error) service.LoyaltyStorage
	}

	makeLoyaltyStorageBalanceOnce := func(t *testing.T, inputUserID uint64, outBalance *model.Balance, outBalanceErr error) service.LoyaltyStorage {
		loyaltyStorage := service.NewMockLoyaltyStorage(t)
		loyaltyStorage.EXPECT().Balance(mock.Anything, inputUserID).Return(outBalance, outBalanceErr).Once()
		return loyaltyStorage
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "StatusOK",
			given: given{
				inputUserID: 1,
				outBalance: &model.Balance{
					ID:        1,
					UserID:    1,
					Current:   19.05,
					Withdrawn: 2019.05,
				},
				outBalanceErr:      nil,
				makeLoyaltyStorage: makeLoyaltyStorageBalanceOnce,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: `{"id":1,"user_id":1,"current":19.05,"withdrawn":2019.05}`,
				contentType:  "application/json",
			},
		},
		/*{
			TODO name: "StatusUnauthorized",
			given: given{
				makeLoyaltyStorage: func(t *testing.T, inputUserID uint64, outBalance *model.Balance, outBalanceErr error) service.LoyaltyStorage {
					return service.NewMockLoyaltyStorage(t)
				},
			},
			want: want{
				httpStatus:   http.StatusUnauthorized,
				responseBody: ``,
			},
		},*/
		{
			name: "StatusInternalServerError",
			given: given{
				inputUserID:        1,
				outBalance:         nil,
				outBalanceErr:      errors.New("some error"),
				makeLoyaltyStorage: makeLoyaltyStorageBalanceOnce,
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

			client := setupServer(t, tt.given.makeLoyaltyStorage(t, tt.given.inputUserID, tt.given.outBalance, tt.given.outBalanceErr))

			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true)
			var result *resty.Response
			var err error
			result, err = req.Get(uriUserBalance)
			assert.NoError(t, err)

			assertResult(t, result, tt.want, tt.given)
		})
	}
}

func Test_PostBalanceWithdraw(t *testing.T) {
	type given struct {
		inputUserID        uint64
		inputOrderNumber   string
		inputSum           float64
		outBalanceErr      error
		makeLoyaltyStorage func(t *testing.T, inputUserID uint64, inputOrderNumber string, inputSum float64, outBalanceErr error) service.LoyaltyStorage
	}

	makeLoyaltyStorageBalanceWithdrawOnce := func(t *testing.T, inputUserID uint64, inputOrderNumber string, inputSum float64, outBalanceErr error) service.LoyaltyStorage {
		loyaltyStorage := service.NewMockLoyaltyStorage(t)
		loyaltyStorage.EXPECT().BalanceWithdraw(mock.Anything, inputUserID, inputOrderNumber, inputSum).Return(outBalanceErr)
		return loyaltyStorage
	}

	makeLoyaltyStorageBalanceWithdrawNever := func(t *testing.T, inputUserID uint64, inputOrderNumber string, inputSum float64, outBalanceErr error) service.LoyaltyStorage {
		return service.NewMockLoyaltyStorage(t)
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "StatusOK",
			given: given{
				inputUserID:        1,
				inputOrderNumber:   "4532015112830366",
				inputSum:           19.05,
				outBalanceErr:      nil,
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawOnce,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		/*{
			TODO name: "StatusUnauthorized",
			given: given{
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawNever,
			},
			want: want{
				httpStatus:   http.StatusUnauthorized,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},*/
		{
			name: "StatusNotFound",
			given: given{
				inputUserID:        1,
				inputOrderNumber:   "4532015112830366",
				inputSum:           19.05,
				outBalanceErr:      perror.ErrOrderNotFound,
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawOnce,
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
				inputUserID:        1,
				inputOrderNumber:   "4532015112830366",
				inputSum:           19.05,
				outBalanceErr:      perror.ErrBalanceInsufficient,
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawOnce,
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
				inputUserID:        1,
				inputOrderNumber:   "234",
				inputSum:           19.05,
				outBalanceErr:      nil,
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawNever,
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
				inputUserID:        1,
				inputOrderNumber:   "4532015112830366",
				inputSum:           0,
				outBalanceErr:      nil,
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawNever,
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
				inputUserID:        1,
				inputOrderNumber:   "4532015112830366",
				inputSum:           -100.25,
				outBalanceErr:      nil,
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawNever,
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
				inputUserID:        1,
				inputOrderNumber:   "4532015112830366",
				inputSum:           100.25,
				outBalanceErr:      errors.New("some error"),
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawOnce,
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

			client := setupServer(t, tt.given.makeLoyaltyStorage(t, tt.given.inputUserID, tt.given.inputOrderNumber, tt.given.inputSum, tt.given.outBalanceErr))

			balanceWithdraw := dto.BalanceWithdraw{
				OrderNumber: tt.given.inputOrderNumber,
				Sum:         tt.given.inputSum,
			}

			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true).SetBody(balanceWithdraw)

			var result *resty.Response
			var err error
			result, err = req.Post(uriUserBalanceWithdraw)
			assert.NoError(t, err)
			assertResult(t, result, tt.want, tt.given)
		})
	}
}

func Test_GetWithdrawals(t *testing.T) {

	type given struct {
		inputUserID        uint64
		outWithdrawals     []model.Withdrawal
		outWithdrawalsErr  error
		makeLoyaltyStorage func(t *testing.T, inputUserID uint64, outWithdrawals []model.Withdrawal, outWithdrawalsErr error) service.LoyaltyStorage
	}

	makeLoyaltyStorageBalanceWithdrawalsOnce := func(t *testing.T, inputUserID uint64, outWithdrawals []model.Withdrawal, outWithdrawalsErr error) service.LoyaltyStorage {
		loyaltyStorage := service.NewMockLoyaltyStorage(t)
		loyaltyStorage.EXPECT().BalanceWithdrawals(mock.Anything, inputUserID).Return(outWithdrawals, outWithdrawalsErr).Once()
		return loyaltyStorage
	}
	/*
	   TODO for StatusUnauthorized
	   	makeLoyaltyStorageBalanceWithdrawalsNever := func(t *testing.T, inputUserID uint64, outWithdrawals []model.Withdrawal, outWithdrawalsErr error) service.LoyaltyStorage {
	   		return service.NewMockLoyaltyStorage(t)
	   	}
	*/
	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "StatusOK",
			given: given{
				inputUserID: 1,
				outWithdrawals: []model.Withdrawal{
					{
						ID:          123,
						UserID:      1,
						Order:       "12345",
						Sum:         123.45,
						ProcessedAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
					},
				},
				outWithdrawalsErr:  nil,
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawalsOnce,
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
				inputUserID:        1,
				outWithdrawals:     nil,
				outWithdrawalsErr:  nil,
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawalsOnce,
			},
			want: want{
				httpStatus:   http.StatusNoContent,
				responseBody: ``,
				contentType:  "application/json",
			},
		},
		/*{
			TODO name: "StatusUnauthorized",
			given: given{
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawalsNever,
			},
			want: want{
				httpStatus:   http.StatusUnauthorized,
				responseBody: ``,
				contentType:  "application/json",
			},
		},*/
		{
			name: "StatusInternalServerError",
			given: given{
				inputUserID:        1,
				outWithdrawals:     nil,
				outWithdrawalsErr:  errors.New("some error"),
				makeLoyaltyStorage: makeLoyaltyStorageBalanceWithdrawalsOnce,
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

			client := setupServer(t, tt.given.makeLoyaltyStorage(t, tt.given.inputUserID, tt.given.outWithdrawals, tt.given.outWithdrawalsErr))
			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true)

			var result *resty.Response
			var err error
			result, err = req.Get(uriUserWithdrawals)
			assert.NoErrorf(t, err, "given: %+v", tt.given)

			assertResult(t, result, tt.want, tt.given)
		})
	}
}
