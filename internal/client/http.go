package client

import (
	"context"
	"net/http"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/go-resty/resty/v2"
)

const urlAccrualOrder = "/api/orders/"

// HTTPAgent - http клиент для взаимодействия с системой начисления баллов
type HTTPAgent struct {
	client *resty.Client
}

func NewClient(appContext *appcontext.AppContext) *HTTPAgent {

	httpClient := resty.New()
	httpClient.SetBaseURL(appContext.ServerConfig.AccrualSystemAddress())

	return &HTTPAgent{
		client: httpClient,
	}
}

// AccrualOrderGet - получение информации о расчёте начислений баллов лояльности
func (a HTTPAgent) AccrualOrderGet(ctx context.Context, orderNumber string) (*http.Response, error) {

	res, err := a.client.R().SetContext(ctx).SetHeader("Accept", "application/json").
		SetDoNotParseResponse(true).
		Get(urlAccrualOrder + orderNumber)
	if err != nil {
		return nil, err
	}

	return res.RawResponse, nil
}
