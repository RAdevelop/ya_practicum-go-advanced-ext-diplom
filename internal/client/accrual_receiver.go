package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/rate_limiter"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/retryer"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/handler"
	"golang.org/x/time/rate"
)

/*
AccrualReceiver - обертка над клиентом взаимодействия с системой начисления баллов.
AccrualReceiver
Выполняет:
  - Получает данные по заказам, для которых надо проверить результат начисления.
  - Асинхронно Обрабатывает ответы от системы начисления, и сохраняет результаты.
*/
type AccrualReceiver struct {
	appContext     *appcontext.AppContext
	loyaltyManager handler.LoyaltyManageable
	pipeline       *Pipeline[dto.Order, dto.AccrualOuter]
}

var errRetry = errors.New("retry")

func NewAccrualReceiver(appContext *appcontext.AppContext, loyaltyManager handler.LoyaltyManageable) *AccrualReceiver {

	maxJobs := uint(10)
	dataBatchSize := uint(100)
	minBackoff := 1 * time.Second
	maxBackoff := 60 * time.Second

	a := &AccrualReceiver{
		appContext:     appContext,
		loyaltyManager: loyaltyManager,
	}

	a.pipeline = NewPipeline(appContext.Logger, maxJobs, dataBatchSize, minBackoff, maxBackoff, a.dataForJobGet(), a.runJob(), a.afterResultCollect())

	return a
}

func (a *AccrualReceiver) Run(ctx context.Context) {
	a.pipeline.Run(ctx)
}

func (a *AccrualReceiver) afterResultCollect() func(ctx context.Context, resList []*dto.AccrualOuter) error {

	return func(ctx context.Context, resList []*dto.AccrualOuter) error {
		//TODO del
		log.Printf("run afterResultCollect with resList: %+v", resList)

		//TODO implement call a.loyaltyManager

		return nil
	}
}

func (a *AccrualReceiver) dataForJobGet() func(ctx context.Context, dataBatchSize uint) ([]dto.Order, error) {

	return func(ctx context.Context, dataBatchSize uint) ([]dto.Order, error) {
		return retryer.RetryLinear(ctx, func(ctx context.Context) ([]dto.Order, error) {

			return a.loyaltyManager.OrdersAwaitingAccrual(ctx, dataBatchSize)

		}, uint(1), new(uint(3)))
	}
}

func (a *AccrualReceiver) runJob() func(ctx context.Context, data dto.Order) (*dto.AccrualOuter, error) {

	initialRateLimit := rate.Limit(5)
	minRateLimit := rate.Limit(0.5)
	maxRateLimit := rate.Limit(2_000)
	burstRateLimit := 5

	client := NewClient(a.appContext)
	rateLimiter := rate_limiter.NewAdaptiveLimiter(initialRateLimit, minRateLimit, maxRateLimit, burstRateLimit)

	return func(ctx context.Context, data dto.Order) (*dto.AccrualOuter, error) {
		for {
			if err := rateLimiter.Wait(ctx); err != nil {
				return nil, err
			}

			accrualOuter, err := a.accrualOrderGet(ctx, client, rateLimiter, data)

			if errors.Is(err, errRetry) {
				continue
			}
			return accrualOuter, err
		}
	}
}

func (a *AccrualReceiver) accrualOrderGet(ctx context.Context, client *HTTPAgent, rateLimiter *rate_limiter.AdaptiveLimiter, data dto.Order) (*dto.AccrualOuter, error) {

	httpResponse, err := retryer.RetryExponential(ctx, func(ctx context.Context) (*http.Response, error) {
		return client.AccrualOrderGet(ctx, data.Number)
	}, uint(2), new(uint(3)))

	if err != nil {
		// сетевая ошибка, таймаут — не связано с rate limit
		return nil, fmt.Errorf("accrual request: %w", err)
	}

	defer func() {
		_ = httpResponse.Body.Close()
	}()

	switch {
	case httpResponse.StatusCode == http.StatusTooManyRequests:
		retryAfter := a.parseRetryAfter(httpResponse)

		rateLimiter.OnRateLimit(ctx, retryAfter)
		return nil, errRetry

	case httpResponse.StatusCode == http.StatusNoContent:
		// заказ не зарегистрирован в системе расчёта
		rateLimiter.OnSuccess()
		return nil, perror.ErrAccrualOrderNotRegistered

	case httpResponse.StatusCode == http.StatusOK:
		accrual, err := a.parseAccrual(httpResponse)
		if err != nil {
			return nil, fmt.Errorf("parse accrual: %w", err)
		}
		rateLimiter.OnSuccess()
		return accrual, nil

	default:
		// 4xx/5xx — не успех, не rate limit
		// не трогаем лимитер
		return nil, fmt.Errorf("unexpected status: %d", httpResponse.StatusCode)
	}
}

func (a *AccrualReceiver) parseRetryAfter(resp *http.Response) time.Duration {
	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		d := time.Until(t)
		if d < 0 {
			return 0
		}
		return d
	}
	return 0
}

func (a *AccrualReceiver) parseAccrual(resp *http.Response) (*dto.AccrualOuter, error) {
	var accrual dto.AccrualOuter
	if err := json.NewDecoder(resp.Body).Decode(&accrual); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &accrual, nil
}
