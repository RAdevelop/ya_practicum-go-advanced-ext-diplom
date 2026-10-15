package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/ratelimiter"
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
	pipeline       *Pipeline[dto.OrderAccrualQueue, dto.AccrualOuter]
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

/*
afterResultCollect - выполняем действия для полученных результатов

В частности, обновляем данные в БД по заказам, балансам...
*/
func (a *AccrualReceiver) afterResultCollect() func(ctx context.Context, accrualResults []*dto.AccrualOuter) error {

	return func(ctx context.Context, accrualResults []*dto.AccrualOuter) error {

		errs := make([]error, 0, len(accrualResults))
		//номера заказов, по которым надо будет удалить или изменить данные в очереди
		orderAccrualQueueUpdate := make([]dto.OrderAccrualQueue, 0, len(accrualResults))

		for _, accrual := range accrualResults {
			if ctx.Err() != nil {
				errs = append(errs, ctx.Err())
				break
			}

			accrualInner := dto.AccrualInner{
				Order:   accrual.Order,
				Accrual: accrual.Accrual,
			}

			accrualInner.Status.ConvertFromOuter(accrual.Status)

			err := a.loyaltyManager.BalanceAccrual(ctx, accrualInner)

			accrualQueue := dto.OrderAccrualQueue{
				OrderNumber: accrualInner.Order,
			}
			switch {
			case err == nil, errors.Is(err, perror.ErrOrderAccrualAlreadyProcessed):

				if !accrual.Status.IsFinal() {
					accrualQueue.NeedUpdate = true
				}

			default:
				errs = append(errs, err)
				continue
			}

			orderAccrualQueueUpdate = append(orderAccrualQueueUpdate, accrualQueue)
		}

		if len(orderAccrualQueueUpdate) > 0 {
			if err := a.loyaltyManager.OrderAccrualQueueUpdate(ctx, orderAccrualQueueUpdate); err != nil {
				errs = append(errs, err)
			}
		}

		return errors.Join(errs...)
	}
}

// dataForJobGet - получаем список заказов, для которых будет проверять результаты начисления
func (a *AccrualReceiver) dataForJobGet() func(ctx context.Context, dataBatchSize uint) ([]dto.OrderAccrualQueue, error) {

	return func(ctx context.Context, dataBatchSize uint) ([]dto.OrderAccrualQueue, error) {
		return retryer.RetryLinear(ctx, func(ctx context.Context) ([]dto.OrderAccrualQueue, error) {

			return a.loyaltyManager.OrderAccrualQueue(ctx, dataBatchSize)

		}, uint(1), new(uint(3)))
	}
}

// runJob - выполняем задачу (запрос к системе начисления для заказа)
func (a *AccrualReceiver) runJob() func(ctx context.Context, data dto.OrderAccrualQueue) (*dto.AccrualOuter, error) {

	initialRateLimit := rate.Limit(5)
	minRateLimit := rate.Limit(0.5)
	maxRateLimit := rate.Limit(2_000)
	burstRateLimit := 5

	client := NewClient(a.appContext)
	rateLimiter := ratelimiter.NewAdaptiveLimiter(initialRateLimit, minRateLimit, maxRateLimit, burstRateLimit)

	return func(ctx context.Context, data dto.OrderAccrualQueue) (*dto.AccrualOuter, error) {
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

// accrualOrderGet - получаем результат начисления по заказу
func (a *AccrualReceiver) accrualOrderGet(ctx context.Context, client *HTTPAgent, rateLimiter *ratelimiter.AdaptiveLimiter, data dto.OrderAccrualQueue) (*dto.AccrualOuter, error) {

	httpResponse, err := retryer.RetryExponential(ctx, func(ctx context.Context) (*http.Response, error) {
		return client.AccrualOrderGet(ctx, data.OrderNumber)
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
		return nil, fmt.Errorf("decode accrual: %w", err)
	}
	return &accrual, nil
}
