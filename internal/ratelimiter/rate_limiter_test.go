package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

/*
ИИ-шка помогала писать
*/

func TestAdaptiveLimiter_SlowStart(t *testing.T) {
	lim := NewAdaptiveLimiter(1, 0.5, 100, 1)
	// successTarget = 20, slowStartThreshold = 100

	for i := 0; i < 20; i++ {
		lim.OnSuccess()
	}
	// после 20 успехов: current=1 < 100 → 1*2 = 2
	assert.Equal(t, rate.Limit(2), lim.Limit())

	for i := 0; i < 20; i++ {
		lim.OnSuccess()
	}
	// 2*2 = 4
	assert.Equal(t, rate.Limit(4), lim.Limit())
}

func TestAdaptiveLimiter_CongestionAvoidance(t *testing.T) {
	lim := NewAdaptiveLimiter(10, 0.5, 100, 1)

	// догоняем до порога (10*2=20, 20*2=40, 40*2=80, 80*2=160→100)
	for i := 0; i < 20*4; i++ {
		lim.OnSuccess()
	}
	// current=100, slowStartThreshold=100 → линейный рост
	// но maxLimit=100 → стоп
	assert.Equal(t, rate.Limit(100), lim.Limit())

	// сбросим через OnRateLimit
	ctx := context.Background()
	lim.OnRateLimit(ctx, 10*time.Millisecond)
	// current упал: 100 → 50, slowStartThreshold=50
	assert.Equal(t, rate.Limit(50), lim.Limit())

	// теперь current=50, ssthresh=50 → линейный рост
	before := lim.Limit()
	for i := 0; i < 20; i++ {
		lim.OnSuccess()
	}
	// stepUp = max(10*0.1, 1) = 1
	assert.Greater(t, lim.Limit(), before)
}

func TestAdaptiveLimiter_OnRateLimit(t *testing.T) {
	lim := NewAdaptiveLimiter(10, 0.5, 100, 1)

	start := time.Now()
	lim.OnRateLimit(context.Background(), 50*time.Millisecond)
	elapsed := time.Since(start)

	// пауза ~50мс
	assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond)

	// лимит упал вдвое: 10 → 5
	assert.Equal(t, rate.Limit(5), lim.Limit())
}

func TestAdaptiveLimiter_MinLimit(t *testing.T) {
	lim := NewAdaptiveLimiter(1, 0.5, 100, 1)

	for i := 0; i < 10; i++ {
		lim.OnRateLimit(context.Background(), time.Millisecond)
	}
	// 1 → 0.5 → 0.5 → ...
	assert.Equal(t, rate.Limit(0.5), lim.Limit())
}

func TestAdaptiveLimiter_PausedIgnoresSecondRateLimit(t *testing.T) {
	lim := NewAdaptiveLimiter(10, 0.5, 100, 1)

	ctx := context.Background()
	go lim.OnRateLimit(ctx, 200*time.Millisecond) // пауза 200мс

	time.Sleep(50 * time.Millisecond) // пока пауза идёт

	// второй 429 — не должен продлить паузу/снизить лимит ещё раз
	start := time.Now()
	lim.OnRateLimit(ctx, 200*time.Millisecond)
	assert.Less(t, time.Since(start), 10*time.Millisecond)

	// ждём окончания первой паузы
	time.Sleep(200 * time.Millisecond)

	// лимит восстановлен до 5, не 2.5
	assert.Equal(t, rate.Limit(5), lim.Limit())
}

func TestAdaptiveLimiter_WaitDuringPause(t *testing.T) {
	lim := NewAdaptiveLimiter(10, 0.5, 100, 1)

	// вычерпать burst: один Wait забирает единственный токен
	assert.NoError(t, lim.Wait(context.Background()))

	// запустить паузу
	go lim.OnRateLimit(context.Background(), 300*time.Millisecond)
	time.Sleep(50 * time.Millisecond) // точно на паузе

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := lim.Wait(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "would exceed context deadline")
}

func TestAdaptiveLimiter_WaitOK(t *testing.T) {
	lim := NewAdaptiveLimiter(100, 0.5, 100, 10) // 100 rps, burst 10

	ctx := context.Background()
	for i := 0; i < 10; i++ {
		assert.NoError(t, lim.Wait(ctx)) // burst хватит
	}
}

func TestAdaptiveLimiter_ResumeAfterPause(t *testing.T) {
	lim := NewAdaptiveLimiter(10, 0.5, 100, 1)

	ctx := context.Background()
	lim.OnRateLimit(ctx, 50*time.Millisecond)

	// после паузы лимит = 5, paused = false
	assert.Equal(t, rate.Limit(5), lim.Limit())

	// Wait должен проходить
	ctx2, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	assert.NoError(t, lim.Wait(ctx2))
}
