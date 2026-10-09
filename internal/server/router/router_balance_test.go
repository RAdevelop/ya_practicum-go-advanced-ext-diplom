package router

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/handler"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetBalance(t *testing.T) {

	type given struct {
		inputCustomerDTO   *dto.Customer
		outBalance         dto.Balance
		outBalanceErr      error
		makeLoyaltyManager func(t *testing.T, inputCustomerDTO *dto.Customer, outBalance dto.Balance, outBalanceErr error) handler.LoyaltyManageable
		authHeaderSet      func(t *testing.T, req *resty.Request, customerDTO *dto.Customer)
	}

	makeLoyaltyManagerBalanceOnce := func(t *testing.T, inputCustomerDTO *dto.Customer, outBalance dto.Balance, outBalanceErr error) handler.LoyaltyManageable {
		loyaltyManager := handler.NewMockLoyaltyManageable(t)
		loyaltyManager.EXPECT().Balance(mock.Anything, inputCustomerDTO).Return(outBalance, outBalanceErr).Once()
		return loyaltyManager
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "StatusOK",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				outBalance: dto.Balance{
					Current:   19.05,
					Withdrawn: 2019.05,
				},
				outBalanceErr:      nil,
				makeLoyaltyManager: makeLoyaltyManagerBalanceOnce,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: `{"current":19.05,"withdrawn":2019.05}`,
				contentType:  "application/json",
			},
		},
		{
			name: "StatusOK",
			given: given{
				inputCustomerDTO:   testCustomerDTO,
				outBalance:         dto.Balance{},
				outBalanceErr:      nil,
				makeLoyaltyManager: makeLoyaltyManagerBalanceOnce,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: `{"current":0,"withdrawn":0}`,
				contentType:  "application/json",
			},
		},
		{
			name: "StatusUnauthorized",
			given: given{
				inputCustomerDTO: nil, // причина StatusUnauthorized
				makeLoyaltyManager: func(t *testing.T, inputCustomerDTO *dto.Customer, outBalance dto.Balance, outBalanceErr error) handler.LoyaltyManageable {
					return handler.NewMockLoyaltyManageable(t)
				},
				authHeaderSet: authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusUnauthorized,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusInternalServerError",
			given: given{
				inputCustomerDTO:   testCustomerDTO,
				outBalance:         dto.Balance{},
				outBalanceErr:      errors.New("some error"),
				makeLoyaltyManager: makeLoyaltyManagerBalanceOnce,
				authHeaderSet:      authHeaderSetCorrect,
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

			client := setupServer(t, tt.given.makeLoyaltyManager(t, tt.given.inputCustomerDTO, tt.given.outBalance, tt.given.outBalanceErr))

			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true)

			tt.given.authHeaderSet(t, req, tt.given.inputCustomerDTO)

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
		inputCustomerDTO   *dto.Customer
		inputWithdraw      dto.BalanceWithdraw
		outBalanceErr      error
		makeLoyaltyManager func(t *testing.T, inputCustomerDTO *dto.Customer, inputWithdraw dto.BalanceWithdraw, outBalanceErr error) handler.LoyaltyManageable
		authHeaderSet      func(t *testing.T, req *resty.Request, customerDTO *dto.Customer)
	}

	makeLoyaltyManagerBalanceWithdrawOnce := func(t *testing.T, inputCustomerDTO *dto.Customer, inputWithdraw dto.BalanceWithdraw, outBalanceErr error) handler.LoyaltyManageable {
		loyaltyManager := handler.NewMockLoyaltyManageable(t)
		loyaltyManager.EXPECT().BalanceWithdraw(mock.Anything, inputCustomerDTO, inputWithdraw).Return(outBalanceErr)
		return loyaltyManager
	}

	makeLoyaltyManagerBalanceWithdrawNever := func(t *testing.T, inputCustomerDTO *dto.Customer, inputWithdraw dto.BalanceWithdraw, outBalanceErr error) handler.LoyaltyManageable {
		return handler.NewMockLoyaltyManageable(t)
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "StatusOK",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				inputWithdraw: dto.BalanceWithdraw{
					OrderNumber: "4532015112830366",
					Sum:         19.05,
				},
				outBalanceErr:      nil,
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawOnce,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusUnauthorized",
			given: given{
				inputCustomerDTO: nil, // причина StatusUnauthorized
				inputWithdraw: dto.BalanceWithdraw{
					OrderNumber: "4532015112830366",
					Sum:         19.05,
				},
				outBalanceErr:      nil,
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawNever,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusUnauthorized,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusUnprocessableEntity order not found by order number",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				inputWithdraw: dto.BalanceWithdraw{
					OrderNumber: "4532015112830366",
					Sum:         19.05,
				},
				outBalanceErr:      perror.ErrOrderNotFound,
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawOnce,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusUnprocessableEntity,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusPaymentRequired",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				inputWithdraw: dto.BalanceWithdraw{
					OrderNumber: "4532015112830366",
					Sum:         19.05,
				},
				outBalanceErr:      perror.ErrBalanceInsufficient,
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawOnce,
				authHeaderSet:      authHeaderSetCorrect,
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
				inputCustomerDTO: testCustomerDTO,
				inputWithdraw: dto.BalanceWithdraw{
					OrderNumber: "234",
					Sum:         19.05,
				},
				outBalanceErr:      nil,
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawNever,
				authHeaderSet:      authHeaderSetCorrect,
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
				inputCustomerDTO: testCustomerDTO,
				inputWithdraw: dto.BalanceWithdraw{
					OrderNumber: "4532015112830366",
					Sum:         0,
				},
				outBalanceErr:      nil,
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawNever,
				authHeaderSet:      authHeaderSetCorrect,
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
				inputCustomerDTO: testCustomerDTO,
				inputWithdraw: dto.BalanceWithdraw{
					OrderNumber: "4532015112830366",
					Sum:         -100.25,
				},
				outBalanceErr:      nil,
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawNever,
				authHeaderSet:      authHeaderSetCorrect,
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
				inputCustomerDTO: testCustomerDTO,
				inputWithdraw: dto.BalanceWithdraw{
					OrderNumber: "4532015112830366",
					Sum:         100.25,
				},
				outBalanceErr:      errors.New("some error"),
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawOnce,
				authHeaderSet:      authHeaderSetCorrect,
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

			client := setupServer(t, tt.given.makeLoyaltyManager(t, tt.given.inputCustomerDTO, tt.given.inputWithdraw, tt.given.outBalanceErr))

			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true).SetBody(tt.given.inputWithdraw)

			tt.given.authHeaderSet(t, req, tt.given.inputCustomerDTO)

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
		inputCustomerDTO   *dto.Customer
		outWithdrawals     []dto.Withdrawal
		outWithdrawalsErr  error
		makeLoyaltyManager func(t *testing.T, inputCustomerDTO *dto.Customer, outWithdrawals []dto.Withdrawal, outWithdrawalsErr error) handler.LoyaltyManageable
		authHeaderSet      func(t *testing.T, req *resty.Request, customerDTO *dto.Customer)
	}

	makeLoyaltyManagerBalanceWithdrawalsOnce := func(t *testing.T, inputCustomerDTO *dto.Customer, outWithdrawals []dto.Withdrawal, outWithdrawalsErr error) handler.LoyaltyManageable {
		loyaltyManager := handler.NewMockLoyaltyManageable(t)
		loyaltyManager.EXPECT().BalanceWithdrawals(mock.Anything, inputCustomerDTO).Return(outWithdrawals, outWithdrawalsErr).Once()
		return loyaltyManager
	}

	makeLoyaltyManagerBalanceWithdrawalsNever := func(t *testing.T, inputCustomerDTO *dto.Customer, outWithdrawals []dto.Withdrawal, outWithdrawalsErr error) handler.LoyaltyManageable {
		return handler.NewMockLoyaltyManageable(t)
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "StatusOK",
			given: given{
				inputCustomerDTO: testCustomerDTO,
				outWithdrawals: []dto.Withdrawal{
					{
						Order:       "12345",
						Sum:         123.45,
						ProcessedAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
					},
				},
				outWithdrawalsErr:  nil,
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawalsOnce,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: `[{"order":"12345","sum":123.45,"processed_at":"2026-09-25T13:12:16Z"}]`,
				contentType:  "application/json",
			},
		},
		{
			name: "StatusNoContent",
			given: given{
				inputCustomerDTO:   testCustomerDTO,
				outWithdrawals:     nil,
				outWithdrawalsErr:  nil,
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawalsOnce,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusNoContent,
				responseBody: ``,
				contentType:  "application/json",
			},
		},
		{
			name: "StatusNoContent",
			given: given{
				inputCustomerDTO:   testCustomerDTO,
				outWithdrawals:     nil,
				outWithdrawalsErr:  perror.ErrWithdrawalNotFound,
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawalsOnce,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusNoContent,
				responseBody: ``,
				contentType:  "application/json",
			},
		},
		{
			name: "StatusUnauthorized",
			given: given{
				inputCustomerDTO:   nil, // причина StatusUnauthorized
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawalsNever,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusUnauthorized,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusInternalServerError",
			given: given{
				inputCustomerDTO:   testCustomerDTO,
				outWithdrawals:     nil,
				outWithdrawalsErr:  errors.New("some error"),
				makeLoyaltyManager: makeLoyaltyManagerBalanceWithdrawalsOnce,
				authHeaderSet:      authHeaderSetCorrect,
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

			client := setupServer(t, tt.given.makeLoyaltyManager(t, tt.given.inputCustomerDTO, tt.given.outWithdrawals, tt.given.outWithdrawalsErr))
			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true)

			tt.given.authHeaderSet(t, req, tt.given.inputCustomerDTO)

			var result *resty.Response
			var err error
			result, err = req.Get(uriUserWithdrawals)
			assert.NoErrorf(t, err, "given: %+v", tt.given)

			assertResult(t, result, tt.want, tt.given)
		})
	}
}
