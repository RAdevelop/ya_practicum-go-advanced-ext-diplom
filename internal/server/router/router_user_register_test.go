package router

import (
	"errors"
	"net/http"
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/jwtoken"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/handler"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UserRegister(t *testing.T) {

	type given struct {
		inputCustomerCredentials any
		outRegisterError         error
		makeLoyaltyManager       func(t *testing.T, inputCustomerCredentials any, outRegisterError error) handler.LoyaltyManageable
	}

	makeLoyaltyManagerUserRegisterNever := func(t *testing.T, inputCustomerCredentials any, outRegisterError error) handler.LoyaltyManageable {
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
				inputCustomerCredentials: dto.CustomerCredentials{
					Login:    testCustomerDTO.Login,
					Password: testCustomerDTO.Login,
				},
				outRegisterError: nil,
				makeLoyaltyManager: func(t *testing.T, inputCustomerCredentials any, outRegisterError error) handler.LoyaltyManageable {
					loyaltyManager := handler.NewMockLoyaltyManageable(t)

					appContext := setupAppContext(t)
					token, err := jwtoken.Generate(*testCustomerDTO, []byte(appContext.ServerConfig.JWTSecret()))
					assert.NoError(t, err)

					customerCredentials, _ := inputCustomerCredentials.(dto.CustomerCredentials)
					loyaltyManager.EXPECT().UserRegister(mock.Anything, customerCredentials, appContext.ServerConfig.JWTSecret()).Return(token, outRegisterError).Once()

					return loyaltyManager
				},
			},
			want: want{
				httpStatus:   http.StatusOK,
				responseBody: ``,
				contentType:  "application/json",
			},
		},
		{
			name: "StatusBadRequest not CustomerCredentials struct",
			given: given{
				inputCustomerCredentials: struct {
					WrongField string
				}{
					WrongField: "",
				},
				outRegisterError:   nil,
				makeLoyaltyManager: makeLoyaltyManagerUserRegisterNever,
			},
			want: want{
				httpStatus:   http.StatusBadRequest,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusBadRequest CustomerCredentials",
			given: given{
				inputCustomerCredentials: dto.CustomerCredentials{
					Login:    "",
					Password: "",
				},
				outRegisterError:   nil,
				makeLoyaltyManager: makeLoyaltyManagerUserRegisterNever,
			},
			want: want{
				httpStatus:   http.StatusBadRequest,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusBadRequest Empty Login",
			given: given{
				inputCustomerCredentials: dto.CustomerCredentials{
					Login:    "",
					Password: testCustomerDTO.Login,
				},
				outRegisterError:   nil,
				makeLoyaltyManager: makeLoyaltyManagerUserRegisterNever,
			},
			want: want{
				httpStatus:   http.StatusBadRequest,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusBadRequest Empty Password",
			given: given{
				inputCustomerCredentials: dto.CustomerCredentials{
					Login:    testCustomerDTO.Login,
					Password: "",
				},
				outRegisterError:   nil,
				makeLoyaltyManager: makeLoyaltyManagerUserRegisterNever,
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
				inputCustomerCredentials: dto.CustomerCredentials{
					Login:    testCustomerDTO.Login,
					Password: testCustomerDTO.Login,
				},
				outRegisterError: perror.ErrCustomerAlreadyExists,
				makeLoyaltyManager: func(t *testing.T, inputCustomerCredentials any, outRegisterError error) handler.LoyaltyManageable {
					loyaltyManager := handler.NewMockLoyaltyManageable(t)

					appContext := setupAppContext(t)

					customerCredentials, _ := inputCustomerCredentials.(dto.CustomerCredentials)
					loyaltyManager.EXPECT().UserRegister(mock.Anything, customerCredentials, appContext.ServerConfig.JWTSecret()).Return("", outRegisterError).Once()

					return loyaltyManager
				},
			},
			want: want{
				httpStatus:   http.StatusConflict,
				responseBody: ``,
				contentType:  "text/plain; charset=utf-8",
			},
		},
		{
			name: "StatusInternalServerError",
			given: given{
				inputCustomerCredentials: dto.CustomerCredentials{
					Login:    testCustomerDTO.Login,
					Password: testCustomerDTO.Login,
				},
				outRegisterError: errors.New("error"),
				makeLoyaltyManager: func(t *testing.T, inputCustomerCredentials any, outRegisterError error) handler.LoyaltyManageable {
					loyaltyStorage := handler.NewMockLoyaltyManageable(t)

					appContext := setupAppContext(t)
					customerCredentials, _ := inputCustomerCredentials.(dto.CustomerCredentials)
					loyaltyStorage.EXPECT().UserRegister(mock.Anything, customerCredentials, appContext.ServerConfig.JWTSecret()).Return("", outRegisterError).Once()
					return loyaltyStorage
				},
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
			client := setupServer(t, tt.given.makeLoyaltyManager(t, tt.given.inputCustomerCredentials, tt.given.outRegisterError))

			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true).
				SetBody(tt.given.inputCustomerCredentials)

			var result *resty.Response
			var err error
			result, err = req.Post(uriUserRegister)
			assert.NoError(t, err)
			assertResult(t, result, tt.want, tt.given)
		})
	}
}
