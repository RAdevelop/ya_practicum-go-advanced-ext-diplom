package router

import (
	"net/http/httptest"
	"testing"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/app_context"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/logger"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/config"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/handler"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/mock"
)

var result *resty.Response
var err error

type given struct {
}

type want struct {
	httpStatus   int
	responseBody string
}

func setupMockLogger(t *testing.T) *logger.MockLogger {
	logMe := logger.NewMockLogger(t)

	//Не знаю как лучше сделать возможное переменное количество параметров для вызова таких методов... :(
	logMe.EXPECT().Info(mock.Anything, mock.Anything, mock.Anything).Maybe()
	logMe.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Maybe()
	logMe.EXPECT().Warn(mock.Anything, mock.Anything, mock.Anything).Maybe()
	logMe.EXPECT().Debug(mock.Anything, mock.Anything, mock.Anything).Maybe()

	return logMe
}

func setupMockConfigServer(t *testing.T) *config.MockProvider {
	cfg := config.NewMockProvider(t)

	cfg.EXPECT().Address().Maybe().Return("localhost:8080")

	return cfg
}

func setupServer(t *testing.T) *resty.Client {
	t.Helper()

	appContext := app_context.New(setupMockLogger(t), setupMockConfigServer(t))
	handlers := handler.New(appContext)
	srv := httptest.NewServer(New(handlers))
	t.Cleanup(srv.Close)

	return resty.New().SetBaseURL(srv.URL)
}
