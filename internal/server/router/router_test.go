package router

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/app_context"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/logger"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/config"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/handler"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/service"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var result *resty.Response
var err error

type want struct {
	httpStatus   int
	responseBody string
	contentType  string
}

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

func setupMockConfigServer(t *testing.T) *config.MockProvider {
	t.Helper()

	cfg := config.NewMockProvider(t)

	cfg.EXPECT().Address().Maybe().Return("localhost:8080")

	return cfg
}

/*
setupServer - создание сервера для фича-тестов
loyaltyStorage - будем создавать моки в зависимости от того, какой кейс будем тестировать
*/
func setupServer(t *testing.T, loyaltyStorage service.LoyaltyStorage) *resty.Client {
	t.Helper()

	appContext := app_context.New(setupMockLogger(t), setupMockConfigServer(t))
	handlers := handler.New(appContext, service.NewLoyaltyManager(loyaltyStorage))
	srv := httptest.NewServer(New(handlers))
	t.Cleanup(srv.Close)

	return resty.New().SetBaseURL(srv.URL)
}

// assertResult - выполняем проверки по результатам выполнения запросов к api
func assertResult(t *testing.T, result *resty.Response, want want, given any) {
	t.Helper()

	assert.Equalf(t, want.httpStatus, result.StatusCode(), "given: %+v", given)
	assert.Equalf(t, want.contentType, result.Header().Get("Content-Type"), "given: %+v", given)

	body, err := io.ReadAll(result.RawResponse.Body)
	assert.NoError(t, err, "given: %+v", given)

	b := strings.TrimSpace(string(body))
	assert.Equalf(t, want.responseBody, b, "given: %+v", given)
	assert.NoErrorf(t, result.RawResponse.Body.Close(), "given: %+v", given)
}
