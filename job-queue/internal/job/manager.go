package job

import (
	"context"
	"fmt"
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
	store Storage
	pool  *worker.Pool
	seq   uint64
}

func NewManager(st Storage, workers, queueSize int) *Manager {
	manager := &Manager{
		store: st,
	}

	manager.pool = worker.NewPool(workers, queueSize, manager.process)

	return manager
}

func (m *Manager) Start(ctx context.Context) {
	m.pool.Start(ctx)
}

func (m *Manager) Submit(payload string) (Job, error) {
	payload = strings.TrimSpace(payload)
	if payload == "" {
		return Job{}, fmt.Errorf("payload is required")
	}

	jobItem := New(m.nextID(), payload)
	m.store.Create(jobItem)

	if err := m.pool.Enqueue(jobItem.ID); err != nil {
		m.store.UpdateStatus(jobItem.ID, StatusFailed, err.Error())
		return Job{}, err
	}

	return jobItem, nil
}

func (m *Manager) Get(id string) (Job, bool) {
	return m.store.Get(id)
}

func (m *Manager) List() []Job {
	return m.store.List()
}

func (m *Manager) process(ctx context.Context, jobID string) {
	if _, ok := m.store.UpdateStatus(jobID, StatusProcessing, ""); !ok {
		return
	}

	select {
	case <-ctx.Done():
		m.store.UpdateStatus(jobID, StatusFailed, "canceled")
		return
	case <-time.After(500 * time.Millisecond):
	}

	jobItem, ok := m.store.Get(jobID)
	if !ok {
		return
	}

	if strings.Contains(strings.ToLower(jobItem.Payload), "fail") {
		m.store.UpdateStatus(jobID, StatusFailed, "simulated failure")
		return
	}

	m.store.UpdateStatus(jobID, StatusDone, "")
}

func (m *Manager) nextID() string {
	seq := atomic.AddUint64(&m.seq, 1)
	return fmt.Sprintf("job-%d-%d", time.Now().UnixNano(), seq)
}
