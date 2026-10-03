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
		makeLoyaltyStorage func(inputUserDTO *dto.User, outOrderList []model.Order, outOrderError error) service.LoyaltyStorage
	}

	makeLoyaltyStorageOrdersOnce := func(inputUserDTO *dto.User, outOrderList []model.Order, outOrderError error) service.LoyaltyStorage {
		loyaltyStorage := service.NewMockLoyaltyStorage(t)
		loyaltyStorage.EXPECT().Orders(mock.Anything, inputUserDTO).Return(outOrderList, outOrderError).Once()

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
				makeLoyaltyStorage: makeLoyaltyStorageOrdersOnce,
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
				makeLoyaltyStorage: makeLoyaltyStorageOrdersOnce,
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
				inputUserDTO: nil,
				makeLoyaltyStorage: func(inputUserDTO *dto.User, outOrderList []model.Order, outOrderError error) service.LoyaltyStorage {
					return service.NewMockLoyaltyStorage(t)
				},
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
				inputUserDTO:       testUserDTO,
				outOrderList:       nil,
				outOrderError:      errors.New("some error"),
				makeLoyaltyStorage: makeLoyaltyStorageOrdersOnce,
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

			client := setupServer(t, tt.given.makeLoyaltyStorage(tt.given.inputUserDTO, tt.given.outOrderList, tt.given.outOrderError), tt.given.inputUserDTO)

			req := client.R().
				SetHeader("Content-Type", "application/json").SetDoNotParseResponse(true)

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
		makeLoyaltyStorage  func(t *testing.T, inputUserDTO *dto.User, inputOrderNumber string, outOrderUploadError error) service.LoyaltyStorage
	}

	makeLoyaltyStorageOrderUploadOnce := func(t *testing.T, inputUserDTO *dto.User, inputOrderNumber string, outOrderUploadError error) service.LoyaltyStorage {
		loyaltyStorage := service.NewMockLoyaltyStorage(t)
		loyaltyStorage.EXPECT().OrderUpload(mock.Anything, inputUserDTO, inputOrderNumber).Return(outOrderUploadError).Once()
		return loyaltyStorage
	}

	makeLoyaltyStorageOrderUploadNever := func(t *testing.T, inputUserDTO *dto.User, inputOrderNumber string, outOrderUploadError error) service.LoyaltyStorage {
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
				inputUserDTO:        testUserDTO,
				inputOrderNumber:    "79927398713",
				outOrderUploadError: perror.ErrOrderAlreadyUploadedByUser,
				makeLoyaltyStorage:  makeLoyaltyStorageOrderUploadOnce,
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
				makeLoyaltyStorage:  makeLoyaltyStorageOrderUploadOnce,
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
				makeLoyaltyStorage:  makeLoyaltyStorageOrderUploadNever,
			},
			want: want{
				httpStatus:   http.StatusBadRequest,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusUnauthorized",
			given: given{
				inputUserDTO:        nil,
				inputOrderNumber:    "4532015112830366",
				outOrderUploadError: nil,
				makeLoyaltyStorage:  makeLoyaltyStorageOrderUploadNever,
			},
			want: want{
				httpStatus:   http.StatusUnauthorized,
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
				makeLoyaltyStorage:  makeLoyaltyStorageOrderUploadOnce,
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
				makeLoyaltyStorage:  makeLoyaltyStorageOrderUploadNever,
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
				makeLoyaltyStorage:  makeLoyaltyStorageOrderUploadOnce,
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

			client := setupServer(t, tt.given.makeLoyaltyStorage(t, tt.given.inputUserDTO, tt.given.inputOrderNumber, tt.given.outOrderUploadError), tt.given.inputUserDTO)

			req := client.R().
				SetHeader("Content-Type", "text/plain").
				SetDoNotParseResponse(true).
				SetBody(tt.given.inputOrderNumber)

			var result *resty.Response
			var err error
			result, err = req.Post(uriUserOrders)
			assert.NoError(t, err)

			assertResult(t, result, tt.want, tt.given)
		})
	}
}
