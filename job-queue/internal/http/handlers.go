package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Godreck/go-pet-projects/job-queue/internal/job"
	"github.com/Godreck/go-pet-projects/job-queue/internal/worker"
)

type Handlers struct {
	manager *job.Manager
	logger  *slog.Logger
}

type createJobRequest struct {
	Payload string `json:"payload"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(manager *job.Manager) http.Handler {
	logger := slog.Default().WithGroup("http")
	h := &Handlers{
		manager: manager,
		logger:  logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.handleHealthz)
	mux.HandleFunc("/jobs", h.handleJobs)
	mux.HandleFunc("/jobs/", h.handleJobByID)

	return mux
}

func (h *Handlers) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handlers) handleJobs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req createJobRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.logger.Warn("invalid request body", "error", err.Error())
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		jobItem, err := h.manager.Submit(req.Payload)
		if err != nil {
			if errors.Is(err, worker.ErrQueueFull) {
				h.logger.Warn("job submit error: queue is full", "payload_len", len(req.Payload))
				writeError(w, http.StatusServiceUnavailable, "queue is full")
				return
			}

			h.logger.Error("job supbit failed", "error", err.Error(), "payload_len", len(req.Payload))
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		h.logger.Info("job submited", "job id", jobItem.ID, "payload_len", len(req.Payload))

		writeJSON(w, http.StatusAccepted, jobItem)
	case http.MethodGet:
		jobs := h.manager.List()
		h.logger.Info("jobs list requested", "count", len(jobs))
		writeJSON(w, http.StatusOK, jobs)
	default:
		h.logger.Warn("method not allowed", "method", r.Method, "path", r.URL.Path)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handlers) handleJobByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Warn("method not allowed", "method", r.Method, "path", r.URL.Path)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/jobs/")
	if id == "" || strings.Contains(id, "/") {
		h.logger.Warn("invalid job id path", "path", r.URL.Path)
		writeError(w, http.StatusNotFound, "job not found")
		return
	}

	jobItem, ok := h.manager.Get(id)
	if !ok {
		h.logger.Warn("job not found", "job_id", id)
		writeError(w, http.StatusNotFound, "job not found")
		return
	}

	h.logger.Info("job requested by id", "job_id", id)
	writeJSON(w, http.StatusOK, jobItem)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, errorResponse{
		Error: message,
	})
}
