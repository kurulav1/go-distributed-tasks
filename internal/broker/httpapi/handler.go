package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/example/distrib-jobs/internal/auth"
	"github.com/example/distrib-jobs/internal/jobs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type Handler struct {
	js       nats.JetStreamContext
	stream   string
	db       *pgxpool.Pool
	auth     *auth.Service
}

type ctxKey string

var userKey ctxKey = "user"

func NewHandler(js nats.JetStreamContext, stream string, dbpool *pgxpool.Pool, a *auth.Service) http.Handler {
	h := &Handler{js: js, stream: stream, db: dbpool, auth: a}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", h.requireAuth(h.enqueue))
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("POST /signup", h.signup)
	mux.HandleFunc("POST /login", h.login)
	return mux
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" || in.Password == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	id, err := h.auth.CreateUser(r.Context(), h.db, in.Email, in.Password)
	if err != nil {
		http.Error(w, "conflict", http.StatusConflict)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	u, err := h.auth.VerifyUser(r.Context(), h.db, in.Email, in.Password)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	tok, err := h.auth.Token(u.ID)
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"token": tok})
}

func (h *Handler) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, err := h.auth.FromRequest(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
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
	uid, _ := r.Context().Value(userKey).(string)
	in.ID = uuid.NewString()
	in.EnqueuedAt = time.Now().UTC()
	in.UserID = uid
	b, err := json.Marshal(in)
	if err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
	subject := "jobs." + in.Type
	if _, err := h.js.Publish(subject, b); err != nil {
		http.Error(w, "broker error", http.StatusServiceUnavailable)
		return
	}
	if err := jobs.Store(r.Context(), h.db, in); err != nil {
		http.Error(w, "db error", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(in)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
