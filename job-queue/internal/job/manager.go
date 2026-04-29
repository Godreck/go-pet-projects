package job

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Godreck/go-pet-projects/job-queue/internal/worker"
)

type Storage interface {
	Create(jobItem Job)
	Get(id string) (Job, bool)
	List() []Job
	UpdateStatus(id string, status Status, errMsg string) (Job, bool)
}

type Manager struct {
	store  Storage
	pool   *worker.Pool
	seq    uint64
	logger *slog.Logger
}

func NewManager(st Storage, workers, queueSize int) *Manager {
	logger := slog.Default().WithGroup("manager")
	manager := &Manager{
		store:  st,
		logger: logger,
	}

	manager.pool = worker.NewPool(workers, queueSize, manager.process)

	return manager
}

func (m *Manager) Start(ctx context.Context) {
	m.pool.SetLogger(m.logger)
	m.pool.Start(ctx)
}

func (m *Manager) Submit(payload string) (Job, error) {
	payload = strings.TrimSpace(payload)
	if payload == "" {
		err := fmt.Errorf("payload is required")
		m.logger.Warn("job submit rejected: payload is empty")
		return Job{}, err
	}

	jobItem := New(m.nextID(), payload)
	m.store.Create(jobItem)

	if err := m.pool.Enqueue(jobItem.ID); err != nil {
		m.logger.Warn("job submit error: queue is full", "job_id", jobItem.ID)
		m.store.UpdateStatus(jobItem.ID, StatusFailed, err.Error())
		return Job{}, err
	}

	m.logger.Info("job submited to pool", "job_id", jobItem.ID, "payload_len", len(payload))

	return jobItem, nil
}

func (m *Manager) Get(id string) (Job, bool) {
	jobItem, ok := m.store.Get(id)
	if !ok {
		m.logger.Warn("job not found in store", "job_id", id)
	}
	return jobItem, ok
}

func (m *Manager) List() []Job {
	jobs := m.store.List()
	m.logger.Debug("jobs list read from store", "count", len(jobs))
	return m.store.List()
}

func (m *Manager) process(ctx context.Context, jobID string) {
	if _, ok := m.store.UpdateStatus(jobID, StatusProcessing, ""); !ok {
		m.logger.Warn("job not found during processing", "job_id", jobID)
		return
	}

	m.logger.Info("job processing started", "job_id", jobID)

	select {
	case <-ctx.Done():
		m.logger.Info("job processing canceled", "job_id", jobID)
		m.store.UpdateStatus(jobID, StatusFailed, "canceled")
		return
	case <-time.After(500 * time.Millisecond):
	}

	jobItem, ok := m.store.Get(jobID)
	if !ok {
		m.logger.Warn("job not found after processing/delay", "job_id", jobID)
		return
	}

	if strings.Contains(strings.ToLower(jobItem.Payload), "fail") {
		m.logger.Warn("job processing simulated failure", "job_id", jobID)
		m.store.UpdateStatus(jobID, StatusFailed, "simulated failure")
		return
	}

	m.store.UpdateStatus(jobID, StatusDone, "")
	m.logger.Info("job processing completed", "job_id", jobID)
}

func (m *Manager) nextID() string {
	seq := atomic.AddUint64(&m.seq, 1)
	return fmt.Sprintf("job-%d-%d", time.Now().UnixNano(), seq)
}
