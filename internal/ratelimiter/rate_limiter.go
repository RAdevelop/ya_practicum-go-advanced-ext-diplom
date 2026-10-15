package ratelimiter

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// AdaptiveLimiter - лимитер подстраивается под работу вызываемого АПИ
type AdaptiveLimiter struct {
	limiter            *rate.Limiter
	mu                 sync.Mutex
	minLimit           rate.Limit // Нижняя граница
	maxLimit           rate.Limit // Верхняя граница
	slowStartThreshold rate.Limit // Порог между slow start и congestion avoidance
	stepUp             rate.Limit // Шаг линейного роста после порога
	successes          int
	successTarget      int
	paused             bool // Лимит на паузе (Retry-After)
}

func NewAdaptiveLimiter(initial rate.Limit, minLimit rate.Limit, maxLimit rate.Limit, burst int) *AdaptiveLimiter {

	if initial <= 0 {
		initial = 1
	}
	if minLimit <= 0 {
		minLimit = 0.5
	}
	if maxLimit < minLimit {
		maxLimit = minLimit
	}
	if initial < minLimit {
		initial = minLimit
	}
	if initial > maxLimit {
		initial = maxLimit
	}

	if burst <= 0 {
		burst = int(initial)
		if burst < 1 {
			burst = 1
		}
	}
	return &AdaptiveLimiter{
		limiter:            rate.NewLimiter(initial, burst),
		minLimit:           minLimit,
		maxLimit:           maxLimit,
		slowStartThreshold: maxLimit, // изначально порог = max
		stepUp:             max(initial*0.1, 1),
		successTarget:      20,
	}
}

// Wait - ждем "разрешения" лимитера на выполнение запроса
func (a *AdaptiveLimiter) Wait(ctx context.Context) error {
	return a.limiter.Wait(ctx)
}

// Limit - получаем текущее значение ограничения
func (a *AdaptiveLimiter) Limit() rate.Limit {
	return a.limiter.Limit()
}

// OnRateLimit - вызывать в случае наступления "события" превышения лимита запросов к АПИ
func (a *AdaptiveLimiter) OnRateLimit(ctx context.Context, retryAfter time.Duration) {

	if retryAfter == 0 {
		retryAfter = time.Second
	}

	var newLimit rate.Limit
	var needPause bool

	func() {
		a.mu.Lock()
		defer a.mu.Unlock()

		// уже на паузе — выходим, пауза уже идёт
		if a.paused {
			return
		}

		current := a.limiter.Limit()
		newLimit = current / 2
		if newLimit < a.minLimit {
			newLimit = a.minLimit
		}
		a.limiter.SetLimit(newLimit)
		a.slowStartThreshold = newLimit // порог падает
		a.successes = 0

		a.paused = true
		a.limiter.SetLimit(0)
		needPause = true
	}()

	// пауза на Retry-After — блокируем всех
	if needPause {
		sleepCtx(ctx, retryAfter)

		a.mu.Lock()
		a.limiter.SetLimit(newLimit)
		a.paused = false
		a.mu.Unlock()
	}
}

// OnSuccess - вызывать в случае успешного запроса к АПИ
func (a *AdaptiveLimiter) OnSuccess() {
	a.mu.Lock()
	defer a.mu.Unlock()

	// пока на паузе — не трогаем лимит
	if a.paused {
		return
	}

	a.successes++
	if a.successes < a.successTarget {
		return
	}
	a.successes = 0

	current := a.limiter.Limit()
	var newLimit rate.Limit

	if current < a.slowStartThreshold {
		// Slow Start: удваиваем
		newLimit = current * 2
	} else {
		// Congestion Avoidance: линейно
		newLimit = current + a.stepUp
	}

	if newLimit > a.maxLimit {
		newLimit = a.maxLimit
	}
	a.limiter.SetLimit(newLimit)
}

// sleepCtx - спим указанное количество времени, или прерываем "сон", если контекст отменен
func sleepCtx(ctx context.Context, retryAfter time.Duration) {

	if retryAfter <= 0 {
		return
	}

	timer := time.NewTimer(retryAfter)
	defer timer.Stop()

	select {
	case <-ctx.Done():

	case <-timer.C:
	}
}
