package client

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// jobResult - структура с результатами работы воркеров
type jobResult[R any] struct {
	Data *R
	Err  error
}

type Pipeline[T any, R any] struct {
	jobs            chan T
	results         chan jobResult[R]
	maxJobs         uint
	dataBatchSize   uint
	logger          Logger
	minBackoff      time.Duration
	maxBackoff      time.Duration
	onJobGet        func(ctx context.Context, dataBatchSize uint) ([]T, error)
	onRunJob        func(ctx context.Context, data T) (r *R, err error)
	onResultCollect func(ctx context.Context, resList []*R) error
}

func NewPipeline[T any, R any](logger Logger, maxJobs uint, dataBatchSize uint,
	minBackoff time.Duration,
	maxBackoff time.Duration,
	onJobGet func(ctx context.Context, dataBatchSize uint) ([]T, error),
	onRunJob func(ctx context.Context, data T) (r *R, err error),
	onResultCollect func(ctx context.Context, resList []*R) error,
) *Pipeline[T, R] {

	if minBackoff <= 0 {
		minBackoff = 1 * time.Second
	}
	if maxBackoff <= 0 {
		maxBackoff = 60 * time.Second
	}

	if maxBackoff < minBackoff {
		cpMaxBackoff := maxBackoff
		maxBackoff = minBackoff
		minBackoff = cpMaxBackoff
	}

	if dataBatchSize == 0 || dataBatchSize > 1_000 {
		dataBatchSize = 1_000
	}

	if maxJobs == 0 || maxJobs > 100 {
		maxJobs = 100
	}

	return &Pipeline[T, R]{
		maxJobs:         maxJobs,
		dataBatchSize:   dataBatchSize,
		jobs:            make(chan T, maxJobs),
		results:         make(chan jobResult[R], maxJobs),
		minBackoff:      minBackoff,
		maxBackoff:      maxBackoff,
		onJobGet:        onJobGet,
		onRunJob:        onRunJob,
		onResultCollect: onResultCollect,
		logger:          logger,
	}
}

// Run - запуск задач
func (p *Pipeline[T, R]) Run(ctx context.Context) {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.jobConsume(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.jobProduce(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.resultCollector(ctx)
	}()
	wg.Wait()
}

// jobProduce - получает данные для задач из некоего источника (onJobGet), и отправляет в работу
func (p *Pipeline[T, R]) jobProduce(ctx context.Context) {
	defer close(p.jobs)

	backoff := p.minBackoff

	for {
		jobList, err := p.onJobGet(ctx, p.dataBatchSize)
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			p.logger.Error("jobProduce: dataForJobGet failed", "error", err, "backoff", backoff)
			if !sleepCtx(ctx, backoff) {
				return
			}
			backoff = min(backoff*2, p.maxBackoff)
			continue
		}

		if len(jobList) == 0 {
			if !sleepCtx(ctx, backoff) {
				return
			}
			backoff = min(backoff*2, p.maxBackoff)
			continue
		}

		// есть данные — работаем на полной скорости
		backoff = p.minBackoff

		for _, job := range jobList {
			select {
			case <-ctx.Done():
				return
			case p.jobs <- job:
			}
		}
	}
}

// sleepCtx- засыпаем на указанное количество времени
func sleepCtx(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func (p *Pipeline[T, R]) jobConsume(ctx context.Context) {

	defer close(p.results)

	var wg sync.WaitGroup

	var i uint
	for i = 0; i < p.maxJobs; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for data := range p.jobs {

				jobRes := p.worker(ctx, data)

				if jobRes.Data == nil && jobRes.Err == nil {
					continue
				}
				p.results <- jobRes
			}
		}()
	}

	wg.Wait()
}

// worker - выполняет p.runJob
func (p *Pipeline[T, R]) worker(ctx context.Context, data T) (jobRes jobResult[R]) {

	defer func() {
		if e := recover(); e != nil {
			jobRes.Err = fmt.Errorf("panicked on task: %v, data: %+v", e, data)
		}
	}()

	jobRes.Data, jobRes.Err = p.onRunJob(ctx, data)

	return
}

// resultCollector - собирает результаты, и выполняем функцию над пачкой результатов
func (p *Pipeline[T, R]) resultCollector(ctx context.Context) {

	batch := make([]*R, 0, p.dataBatchSize)
	for result := range p.results {

		if result.Err != nil {
			/*
				По-хорошему, тут нужно бизнес-правило, что делать с заказом в такой ситуации. Например:
				- если в течении какого-то времени по заказу получаем ошибки от системы начисления, то:
					- удалять из очереди такие заказы по достижении ошибок в количестве N штук, и перекладывать в DLQ очередь для отдельного разбора/анализа
					- до достижения ошибок в количестве N штук, обновляем заказ в очереди - откладываем по нему следующую выборку
				- в этом случае, надо данные для канала results (type jobResult[R any] struct {Data *R Err error}) дополнить входными данными:
				type jobResult[R any] struct {inputData *In Data *R Err error}
					- тогда по inputData (dto.OrderAccrualQueue) можно будет обработать такую ситуацию, зная номер заказа.
					- в этом случае, текущая обработка if result.Err != nil - будет не нужна. Точнее, она перейдет в реализацию onResultCollect() вызывающей стороны.
			*/
			p.logger.Error("resultCollector: result failed", "error", result.Err)
			continue
		}

		batch = append(batch, result.Data)

		if uint(len(batch)) >= p.dataBatchSize {
			err := p.onResultCollect(ctx, batch)
			if err != nil {
				p.logger.Error("resultCollector: result failed", "error", err)
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		err := p.onResultCollect(ctx, batch)
		if err != nil {
			p.logger.Error("resultCollector: result failed", "err", err)
		}
	}
}
