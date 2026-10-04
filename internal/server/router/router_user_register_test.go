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
		inputUserCredentials any
		outRegisterError     error
		makeLoyaltyManager   func(t *testing.T, inputUserCredentials any, outRegisterError error) handler.LoyaltyManageable
	}

	makeLoyaltyManagerUserRegisterNever := func(t *testing.T, inputUserCredentials any, outRegisterError error) handler.LoyaltyManageable {
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
				inputUserCredentials: dto.UserCredentials{
					Login:    testUserDTO.Login,
					Password: testUserDTO.Login,
				},
				outRegisterError: nil,
				makeLoyaltyManager: func(t *testing.T, inputUserCredentials any, outRegisterError error) handler.LoyaltyManageable {
					loyaltyManager := handler.NewMockLoyaltyManageable(t)

					appContext := setupAppContext(t)
					token, err := jwtoken.Generate(*testUserDTO, []byte(appContext.ServerConfig.JWTSecret()))
					assert.NoError(t, err)

					userCredentials, _ := inputUserCredentials.(dto.UserCredentials)
					loyaltyManager.EXPECT().UserRegister(mock.Anything, userCredentials, appContext.ServerConfig.JWTSecret()).Return(token, outRegisterError).Once()

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
			name: "StatusBadRequest not UserCredentials struct",
			given: given{
				inputUserCredentials: struct {
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
			name: "StatusBadRequest UserCredentials",
			given: given{
				inputUserCredentials: dto.UserCredentials{
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
				inputUserCredentials: dto.UserCredentials{
					Login:    "",
					Password: testUserDTO.Login,
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
				inputUserCredentials: dto.UserCredentials{
					Login:    testUserDTO.Login,
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
				inputUserCredentials: dto.UserCredentials{
					Login:    testUserDTO.Login,
					Password: testUserDTO.Login,
				},
				outRegisterError: perror.ErrUserAlreadyExists,
				makeLoyaltyManager: func(t *testing.T, inputUserCredentials any, outRegisterError error) handler.LoyaltyManageable {
					loyaltyManager := handler.NewMockLoyaltyManageable(t)

					appContext := setupAppContext(t)

					userCredentials, _ := inputUserCredentials.(dto.UserCredentials)
					loyaltyManager.EXPECT().UserRegister(mock.Anything, userCredentials, appContext.ServerConfig.JWTSecret()).Return("", outRegisterError).Once()

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
				inputUserCredentials: dto.UserCredentials{
					Login:    testUserDTO.Login,
					Password: testUserDTO.Login,
				},
				outRegisterError: errors.New("error"),
				makeLoyaltyManager: func(t *testing.T, inputUserCredentials any, outRegisterError error) handler.LoyaltyManageable {
					loyaltyStorage := handler.NewMockLoyaltyManageable(t)

					appContext := setupAppContext(t)
					userCredentials, _ := inputUserCredentials.(dto.UserCredentials)
					loyaltyStorage.EXPECT().UserRegister(mock.Anything, userCredentials, appContext.ServerConfig.JWTSecret()).Return("", outRegisterError).Once()
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
			client := setupServer(t, tt.given.makeLoyaltyManager(t, tt.given.inputUserCredentials, tt.given.outRegisterError))

			req := client.R().
				SetHeader("Content-Type", "application/json").
				SetDoNotParseResponse(true).
				SetBody(tt.given.inputUserCredentials)

			var result *resty.Response
			var err error
			result, err = req.Post(uriUserRegister)
			assert.NoError(t, err)
			assertResult(t, result, tt.want, tt.given)
		})
	}
}
