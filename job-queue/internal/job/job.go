package job

import "time"

type Status string

const (
	StatusQueued     Status = "queued"
	StatusProcessing Status = "processing"
	StatusDone       Status = "done"
	StatusFailed     Status = "failed"
)

type Job struct {
	ID        string    `json:"id"`
	Payload   string    `json:"payload"`
	Status    Status    `json:"status"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(id, payload string) Job {
	now := time.Now().UTC()

	return Job{
		ID:        id,
		Payload:   payload,
		Status:    StatusQueued,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
