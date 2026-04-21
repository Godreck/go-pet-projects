package worker

import (
	"context"
	"errors"
	"sync"
)

var ErrQueueFull = errors.New("queue is full")

type Processor func(ctx context.Context, jobID string)

type Pool struct {
	workerCount int
	queue       chan string
	processor   Processor
	wg          sync.WaitGroup
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

	return &Pool{
		workerCount: workerCount,
		queue:       make(chan string, queueSize),
		processor:   processor,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case jobID := <-p.queue:
					p.processor(ctx, jobID)
				}
			}
		}()
	}
}

func (p *Pool) Enqueue(jobID string) error {
	select {
	case p.queue <- jobID:
		return nil
	default:
		return ErrQueueFull
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
