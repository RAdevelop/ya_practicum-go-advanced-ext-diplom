package router

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/handler"
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetOrders(t *testing.T) {

	type given struct {
		inputCustomerDTO   *dto.Customer
		outOrderList       []model.Order
		outOrderError      error
		makeLoyaltyManager func(t *testing.T, inputCustomerDTO *dto.Customer, outOrderList []model.Order, outOrderError error) handler.LoyaltyManageable
		authHeaderSet      func(t *testing.T, req *resty.Request, customerDTO *dto.Customer)
	}

	makeLoyaltyManagerOrdersOnce := func(t *testing.T, inputCustomerDTO *dto.Customer, outOrderList []model.Order, outOrderError error) handler.LoyaltyManageable {
		loyaltyManager := handler.NewMockLoyaltyManageable(t)
		loyaltyManager.EXPECT().Orders(mock.Anything, inputCustomerDTO).Return(outOrderList, outOrderError).Once()

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
				outOrderList: []model.Order{
					{
						ID:         1,
						Number:     "12345",
						CustomerID: 1,
						Status:     statusInner.AccrualNew,
						UploadedAt: time.Date(2026, 9, 25, 13, 12, 16, 0, time.UTC),
					},
				},
				outOrderError:      nil,
				makeLoyaltyManager: makeLoyaltyManagerOrdersOnce,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: `[{"id":1,"number":"12345","customer_id":1,"status":"NEW","uploaded_at":"2026-09-25T13:12:16Z"}]`,
				contentType:  "application/json",
			},
		},
		{
			name: "StatusNoContent",
			given: given{
				inputCustomerDTO:   testCustomerDTO,
				outOrderList:       nil,
				outOrderError:      nil,
				makeLoyaltyManager: makeLoyaltyManagerOrdersOnce,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusNoContent,
				responseBody: ``,
				contentType:  "application/json",
			},
		},
		{
			name: "StatusInternalServerError",
			given: given{
				inputCustomerDTO:   testCustomerDTO,
				outOrderList:       nil,
				outOrderError:      errors.New("some error"),
				makeLoyaltyManager: makeLoyaltyManagerOrdersOnce,
				authHeaderSet:      authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusInternalServerError,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusUnauthorized",
			given: given{
				inputCustomerDTO: nil, // причина StatusUnauthorized
				makeLoyaltyManager: func(t *testing.T, inputCustomerDTO *dto.Customer, outOrderList []model.Order, outOrderError error) handler.LoyaltyManageable {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			client := setupServer(t, tt.given.makeLoyaltyManager(t, tt.given.inputCustomerDTO, tt.given.outOrderList, tt.given.outOrderError))

			req := client.R().
				SetHeader("Content-Type", "application/json").SetDoNotParseResponse(true)

			tt.given.authHeaderSet(t, req, tt.given.inputCustomerDTO)

			var result *resty.Response
			var err error
			result, err = req.Get(uriUserOrders)
			assert.NoError(t, err)

			assertResult(t, result, tt.want, tt.given)
		})
	}
}

func Test_PostOrders(t *testing.T) {

	type given struct {
		inputCustomerDTO    *dto.Customer
		inputOrderNumber    string
		outOrderUploadError error
		makeLoyaltyManager  func(t *testing.T, inputCustomerDTO *dto.Customer, inputOrderNumber string, outOrderUploadError error) handler.LoyaltyManageable
		authHeaderSet       func(t *testing.T, req *resty.Request, customerDTO *dto.Customer)
	}

	makeLoyaltyManagerOrderUploadOnce := func(t *testing.T, inputCustomerDTO *dto.Customer, inputOrderNumber string, outOrderUploadError error) handler.LoyaltyManageable {
		loyaltyManager := handler.NewMockLoyaltyManageable(t)
		loyaltyManager.EXPECT().OrderUpload(mock.Anything, inputCustomerDTO, inputOrderNumber).Return(outOrderUploadError).Once()
		return loyaltyManager
	}

	makeLoyaltyManagerOrderUploadNever := func(t *testing.T, inputCustomerDTO *dto.Customer, inputOrderNumber string, outOrderUploadError error) handler.LoyaltyManageable {
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
				inputCustomerDTO:    testCustomerDTO,
				inputOrderNumber:    "79927398713",
				outOrderUploadError: perror.ErrOrderAlreadyUploadedByCustomer,
				makeLoyaltyManager:  makeLoyaltyManagerOrderUploadOnce,
				authHeaderSet:       authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusAccepted",
			given: given{
				inputCustomerDTO:    testCustomerDTO,
				inputOrderNumber:    "4532015112830366",
				outOrderUploadError: nil,
				makeLoyaltyManager:  makeLoyaltyManagerOrderUploadOnce,
				authHeaderSet:       authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusAccepted,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusBadRequest",
			given: given{
				inputCustomerDTO:    testCustomerDTO,
				inputOrderNumber:    "",
				outOrderUploadError: nil,
				makeLoyaltyManager:  makeLoyaltyManagerOrderUploadNever,
				authHeaderSet:       authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusBadRequest,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusConflict",
			given: given{
				inputCustomerDTO:    testCustomerDTO,
				inputOrderNumber:    "4532015112830366",
				outOrderUploadError: perror.ErrOrderAlreadyUploadedByOther,
				makeLoyaltyManager:  makeLoyaltyManagerOrderUploadOnce,
				authHeaderSet:       authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusConflict,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusUnprocessableEntity",
			given: given{
				inputCustomerDTO:    testCustomerDTO,
				inputOrderNumber:    "12345string",
				outOrderUploadError: nil,
				makeLoyaltyManager:  makeLoyaltyManagerOrderUploadNever,
				authHeaderSet:       authHeaderSetCorrect,
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
				inputCustomerDTO:    testCustomerDTO,
				inputOrderNumber:    "4532015112830366",
				outOrderUploadError: errors.New("some error"),
				makeLoyaltyManager:  makeLoyaltyManagerOrderUploadOnce,
				authHeaderSet:       authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusInternalServerError,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusUnauthorized",
			given: given{
				inputCustomerDTO:    nil, // причина StatusUnauthorized
				inputOrderNumber:    "4532015112830366",
				outOrderUploadError: nil,
				makeLoyaltyManager:  makeLoyaltyManagerOrderUploadNever,
				authHeaderSet:       authHeaderSetCorrect,
			},
			want: want{
				httpStatus:   http.StatusUnauthorized,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			client := setupServer(t, tt.given.makeLoyaltyManager(t, tt.given.inputCustomerDTO, tt.given.inputOrderNumber, tt.given.outOrderUploadError))

			req := client.R().
				SetHeader("Content-Type", "text/plain").
				SetDoNotParseResponse(true).
				SetBody(tt.given.inputOrderNumber)

			tt.given.authHeaderSet(t, req, tt.given.inputCustomerDTO)

			var result *resty.Response
			var err error
			result, err = req.Post(uriUserOrders)
			assert.NoError(t, err)

			assertResult(t, result, tt.want, tt.given)
		})
	}
}
