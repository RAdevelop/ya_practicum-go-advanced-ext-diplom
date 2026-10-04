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
		inputUserDTO       *dto.User
		outOrderList       []model.Order
		outOrderError      error
		makeLoyaltyManager func(t *testing.T, inputUserDTO *dto.User, outOrderList []model.Order, outOrderError error) handler.LoyaltyManageable
		authHeaderSet      func(t *testing.T, req *resty.Request, userDTO *dto.User)
	}

	makeLoyaltyManagerOrdersOnce := func(t *testing.T, inputUserDTO *dto.User, outOrderList []model.Order, outOrderError error) handler.LoyaltyManageable {
		loyaltyManager := handler.NewMockLoyaltyManageable(t)
		loyaltyManager.EXPECT().Orders(mock.Anything, inputUserDTO).Return(outOrderList, outOrderError).Once()

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
				inputUserDTO: testUserDTO,
				outOrderList: []model.Order{
					{
						Number:     "12345",
						UserID:     1,
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
				responseBody: `[{"number":"12345","user_id":1,"status":"NEW","uploaded_at":"2026-09-25T13:12:16Z"}]`,
				contentType:  "application/json",
			},
		},
		{
			name: "StatusNoContent",
			given: given{
				inputUserDTO:       testUserDTO,
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
				inputUserDTO:       testUserDTO,
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
				inputUserDTO: nil, // причина StatusUnauthorized
				makeLoyaltyManager: func(t *testing.T, inputUserDTO *dto.User, outOrderList []model.Order, outOrderError error) handler.LoyaltyManageable {
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

			client := setupServer(t, tt.given.makeLoyaltyManager(t, tt.given.inputUserDTO, tt.given.outOrderList, tt.given.outOrderError))

			req := client.R().
				SetHeader("Content-Type", "application/json").SetDoNotParseResponse(true)

			tt.given.authHeaderSet(t, req, tt.given.inputUserDTO)

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
		inputUserDTO        *dto.User
		inputOrderNumber    string
		outOrderUploadError error
		makeLoyaltyManager  func(t *testing.T, inputUserDTO *dto.User, inputOrderNumber string, outOrderUploadError error) handler.LoyaltyManageable
		authHeaderSet       func(t *testing.T, req *resty.Request, userDTO *dto.User)
	}

	makeLoyaltyManagerOrderUploadOnce := func(t *testing.T, inputUserDTO *dto.User, inputOrderNumber string, outOrderUploadError error) handler.LoyaltyManageable {
		loyaltyManager := handler.NewMockLoyaltyManageable(t)
		loyaltyManager.EXPECT().OrderUpload(mock.Anything, inputUserDTO, inputOrderNumber).Return(outOrderUploadError).Once()
		return loyaltyManager
	}

	makeLoyaltyManagerOrderUploadNever := func(t *testing.T, inputUserDTO *dto.User, inputOrderNumber string, outOrderUploadError error) handler.LoyaltyManageable {
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
				inputUserDTO:        testUserDTO,
				inputOrderNumber:    "79927398713",
				outOrderUploadError: perror.ErrOrderAlreadyUploadedByUser,
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
				inputUserDTO:        testUserDTO,
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
				inputUserDTO:        testUserDTO,
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
				inputUserDTO:        testUserDTO,
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
				inputUserDTO:        testUserDTO,
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
				inputUserDTO:        testUserDTO,
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
				inputUserDTO:        nil, // причина StatusUnauthorized
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

			client := setupServer(t, tt.given.makeLoyaltyManager(t, tt.given.inputUserDTO, tt.given.inputOrderNumber, tt.given.outOrderUploadError))

			req := client.R().
				SetHeader("Content-Type", "text/plain").
				SetDoNotParseResponse(true).
				SetBody(tt.given.inputOrderNumber)

			tt.given.authHeaderSet(t, req, tt.given.inputUserDTO)

			var result *resty.Response
			var err error
			result, err = req.Post(uriUserOrders)
			assert.NoError(t, err)

			assertResult(t, result, tt.want, tt.given)
		})
	}
}
