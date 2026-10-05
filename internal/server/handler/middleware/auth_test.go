package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/jwtoken"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/logger"
	dbConfig "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/repository/database/config"
	serverConfig "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// testCustomerDTO - эмуляция пользователя
var testCustomerDTO = &dto.Customer{
	ID:    1,
	Login: "TestLogin",
}

const testJWTSecret = "lPD7WBZ/MCBKK0aEqgzSqfIQSqAGB7VhIfjsZwuXLJE="

func setupMockLogger(t *testing.T) *logger.MockLogger {
	t.Helper()

	logMe := logger.NewMockLogger(t)

	//Не знаю как лучше сделать возможное переменное количество параметров для вызова таких методов... :(
	logMe.EXPECT().Info(mock.Anything, mock.Anything, mock.Anything).Maybe()
	logMe.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Maybe()
	logMe.EXPECT().Warn(mock.Anything, mock.Anything, mock.Anything).Maybe()
	logMe.EXPECT().Debug(mock.Anything, mock.Anything, mock.Anything).Maybe()

	return logMe
}

func setupMockConfigServer(t *testing.T) *serverConfig.MockProvider {
	t.Helper()

	cfg := serverConfig.NewMockProvider(t)

	cfg.EXPECT().Address().Maybe().Return("localhost:8080")
	cfg.EXPECT().AccrualSystemAddress().Maybe().Return("localhost:8081")
	cfg.EXPECT().JWTSecret().Maybe().Return(testJWTSecret)

	return cfg
}

func setupMockConfigDB(t *testing.T) *dbConfig.MockProvider {
	t.Helper()

	cfg := dbConfig.NewMockProvider(t)

	cfg.EXPECT().DSN().Maybe().Return("db_address")
	cfg.EXPECT().MaxConns().Maybe().Return(25)
	cfg.EXPECT().MinConns().Maybe().Return(5)
	cfg.EXPECT().MaxConnLifetime().Maybe().Return("1h")
	cfg.EXPECT().MaxConnIdleTime().Maybe().Return("4m")

	return cfg
}

func setupAppContext(t *testing.T) *appcontext.AppContext {
	t.Helper()

	return appcontext.New(setupMockLogger(t), setupMockConfigServer(t), setupMockConfigDB(t))
}

func TestAuth(t *testing.T) {
	type given struct {
		header   string
		customer *dto.Customer // Если задан — генерируем токен
	}

	tests := []struct {
		name         string
		given        given
		wantStatus   int
		wantNextCall bool
	}{
		{
			name:         "no Authorization header",
			given:        given{header: ""},
			wantStatus:   http.StatusUnauthorized,
			wantNextCall: false,
		},
		{
			name:         "wrong scheme",
			given:        given{header: "Basic dXNlcjpwYXNz"},
			wantStatus:   http.StatusUnauthorized,
			wantNextCall: false,
		},
		{
			name:         "invalid token",
			given:        given{header: "Bearer invalid-token"},
			wantStatus:   http.StatusUnauthorized,
			wantNextCall: false,
		},
		{
			name: "valid token",
			given: given{
				customer: testCustomerDTO,
			},
			wantStatus:   http.StatusOK,
			wantNextCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appContext := setupAppContext(t)

			// Флаг «next вызван» и захваченный customerDTO.
			var nextCalled bool
			var customerDTO *dto.Customer

			// фейк-хендлер помощник проверки - дошел до него вызов или нет
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				customerDTO, _ = CustomerGetFromCtx(r.Context())
				w.WriteHeader(http.StatusOK)
			})

			authHandler := Auth(appContext, next)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.given.customer != nil {
				token, err := jwtoken.Generate(*tt.given.customer, []byte(testJWTSecret))
				assert.NoErrorf(t, err, "given: %+v", tt.given)
				req.Header.Set("Authorization", "Bearer "+token)
			} else if tt.given.header != "" {
				req.Header.Set("Authorization", tt.given.header)
			}

			rec := httptest.NewRecorder()
			authHandler.ServeHTTP(rec, req)

			assert.Equalf(t, tt.wantStatus, rec.Code, "given: %+v", tt.given)

			if tt.wantNextCall {
				assert.Truef(t, nextCalled, "next должен быть вызван, given: %+v", tt.given)
				assert.NotNilf(t, customerDTO, "given: %+v", tt.given)
				assert.Equalf(t, tt.given.customer, customerDTO, "given: %+v", tt.given)
			} else {
				assert.Falsef(t, nextCalled, "next НЕ должен быть вызван, given: %+v", tt.given)
			}
		})
	}
}

func Test_CustomerDTOGetFromCtx(t *testing.T) {
	type given struct {
		putData any
	}

	type want struct {
		customerDTO *dto.Customer
		expectedOk  bool
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "get dto.Customer",
			given: given{
				putData: testCustomerDTO,
			},
			want: want{
				customerDTO: testCustomerDTO,
				expectedOk:  true,
			},
		},
		{
			name: "do not get dto.Customer",
			given: given{
				putData: "some data",
			},
			want: want{
				customerDTO: nil,
				expectedOk:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.WithValue(t.Context(), keyCustomerDTO, tt.given.putData)

			customerDTO, ok := CustomerGetFromCtx(ctx)
			assert.Equalf(t, tt.want.expectedOk, ok, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.customerDTO, customerDTO, "given: %+v", tt.given)
		})
	}
}

func Test_CustomerDTOPutToCtx(t *testing.T) {
	type given struct {
		customerDTO *dto.Customer
	}

	type want struct {
		customerDTO *dto.Customer
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "put not empty dto.Customer",
			given: given{
				customerDTO: testCustomerDTO,
			},
			want: want{
				customerDTO: testCustomerDTO,
			},
		},
		{
			name: "put empty dto.Customer",
			given: given{
				customerDTO: &dto.Customer{},
			},
			want: want{
				customerDTO: &dto.Customer{},
			},
		},
		{
			name: "put nil dto.Customer",
			given: given{
				customerDTO: nil,
			},
			want: want{
				customerDTO: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := CustomerPutToCtx(t.Context(), tt.given.customerDTO)
			customerDTO, ok := CustomerGetFromCtx(ctx)
			assert.Equalf(t, true, ok, "given: %+v", tt.given)
			assert.Equalf(t, tt.want.customerDTO, customerDTO, "given: %+v", tt.given)
		})
	}
}
