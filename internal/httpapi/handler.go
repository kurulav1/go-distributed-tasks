package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/example/distrib-jobs/internal/jobs"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type Handler struct {
	js     nats.JetStreamContext
	stream string
}

func NewHandler(js nats.JetStreamContext, stream string) http.Handler {
	h := &Handler{js: js, stream: stream}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", h.enqueue)
	mux.HandleFunc("GET /healthz", h.health)
	return mux
}

func (h *Handler) enqueue(w http.ResponseWriter, r *http.Request) {
	var in jobs.Job
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if in.Type == "" || len(in.Payload) == 0 {
		http.Error(w, "invalid job", http.StatusBadRequest)
		return
	}
	in.ID = uuid.NewString()
	in.EnqueuedAt = time.Now().UTC()
	b, err := json.Marshal(in)
	if err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
	subject := "jobs." + in.Type
	_, err = h.js.Publish(subject, b)
	if err != nil {
		http.Error(w, "broker error", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(in)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
