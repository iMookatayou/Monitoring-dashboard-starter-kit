package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	// >>> เพิ่ม import นี้ <<<
	"github.com/iMookatayou/monitoring-dashboard-starter-kit/monitoring-backend/internal/repository"
)

type Repository interface {
	InsertMetric(ctx context.Context, m map[string]any) error

	// >>> แก้บรรทัดนี้ให้คืน []repository.Metric <<<
	QueryMetrics(ctx context.Context, service, name string, from, to time.Time, limit int) ([]repository.Metric, error)

	Ping(ctx context.Context) error
}

type HTTP struct{ Repo Repository }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func readJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	defer r.Body.Close()
	b, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "read body"})
		return false
	}
	if err := json.Unmarshal(b, dst); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return false
	}
	return true
}

// GET /healthz
func (h HTTP) Healthz(w http.ResponseWriter, r *http.Request) {
	if err := h.Repo.Ping(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"status": "db_error", "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /metrics
func (h HTTP) Ingest(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any
	if !readJSON(w, r, &payload) {
		return
	}
	if _, ok := payload["observed_at"]; !ok {
		payload["observed_at"] = time.Now().UTC()
	}
	if err := h.Repo.InsertMetric(r.Context(), payload); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "ok"})
}

// GET /metrics?service=...&name=...&from=RFC3339&to=RFC3339
func (h HTTP) Query(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	service := q.Get("service")
	name := q.Get("name")
	from, to := time.Now().Add(-1*time.Hour), time.Now()
	if s := q.Get("from"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			from = t
		}
	}
	if s := q.Get("to"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			to = t
		}
	}
	rows, err := h.Repo.QueryMetrics(r.Context(), service, name, from, to, 500)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, rows)
}
