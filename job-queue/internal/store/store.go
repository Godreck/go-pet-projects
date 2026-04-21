package store

import (
	"sync"
	"time"

	"github.com/Godreck/go-pet-projects/job-queue/internal/job"
)

type Store struct {
	mu   sync.RWMutex
	jobs map[string]job.Job
}

func New() *Store {
	return &Store{
		jobs: make(map[string]job.Job),
	}
}

func (s *Store) Create(jobItem job.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs[jobItem.ID] = jobItem
}

func (s *Store) Get(id string) (job.Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobItem, ok := s.jobs[id]
	return jobItem, ok
}

func (s *Store) List() []job.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]job.Job, 0, len(s.jobs))
	for _, item := range s.jobs {
		result = append(result, item)
	}

	return result
}

func (s *Store) UpdateStatus(id string, status job.Status, errMsg string) (job.Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobItem, ok := s.jobs[id]
	if !ok {
		return job.Job{}, false
	}

	jobItem.Status = status
	jobItem.Error = errMsg
	jobItem.UpdatedAt = time.Now().UTC()
	s.jobs[id] = jobItem

	return jobItem, true
}
