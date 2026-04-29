package worker

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

var ErrQueueFull = errors.New("queue is full")

type Processor func(ctx context.Context, jobID string)

type Pool struct {
	workerCount int
	queue       chan string
	processor   Processor
	wg          sync.WaitGroup
	logger      *slog.Logger
}

func NewPool(workerCount, queueSize int, processor Processor) *Pool {
	if workerCount < 1 {
		workerCount = 1
	}

	if queueSize < 1 {
		queueSize = 1
	}

	if processor == nil {
		processor = func(context.Context, string) {}
	}

	logger := slog.Default().WithGroup("worker_pool")
	return &Pool{
		workerCount: workerCount,
		queue:       make(chan string, queueSize),
		processor:   processor,
		logger:      logger,
	}
}

func (p *Pool) SetLogger(logger *slog.Logger) {
	p.logger = logger.WithGroup("worker_pool")
}

func (p *Pool) Start(ctx context.Context) {
	p.logger.Info("worker pool started", "workers", p.workerCount, "queue_size", cap(p.queue))

	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			for {
				select {
				case <-ctx.Done():
					p.logger.Info("worker exiting due to context done")
					return
				case jobID := <-p.queue:
					p.logger.Debug("worker got job", "job_id", jobID)
					p.processor(ctx, jobID)
				}
			}
		}()
	}
}

func (p *Pool) Enqueue(jobID string) error {
	select {
	case p.queue <- jobID:
		p.logger.Debug("job enqueued", "job_id", jobID)
		return nil
	default:
		p.logger.Warn("job enqueue faild: queue full", "job_id", jobID)
		return ErrQueueFull
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
