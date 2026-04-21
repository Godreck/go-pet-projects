package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Godreck/go-pet-projects/job-queue/internal/job"
	"github.com/Godreck/go-pet-projects/job-queue/internal/worker"
)

type Handlers struct {
	manager *job.Manager
}

type createJobRequest struct {
	Payload string `json:"payload"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(manager *job.Manager) http.Handler {
	h := &Handlers{
		manager: manager,
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
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		jobItem, err := h.manager.Submit(req.Payload)
		if err != nil {
			if errors.Is(err, worker.ErrQueueFull) {
				writeError(w, http.StatusServiceUnavailable, "queue is full")
				return
			}

			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusAccepted, jobItem)
	case http.MethodGet:
		writeJSON(w, http.StatusOK, h.manager.List())
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handlers) handleJobByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/jobs/")
	if id == "" || strings.Contains(id, "/") {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}

	jobItem, ok := h.manager.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}

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
