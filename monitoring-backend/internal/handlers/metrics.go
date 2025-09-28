package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	InsertMetric(ctx context.Context, m map[string]any) error
	QueryMetrics(ctx context.Context, service, name string, from, to time.Time, limit int) ([]struct{
		ID int64 `db:"id"`
		Service string `db:"service"`
		Name string `db:"name"`
		Value float64 `db:"value"`
		LabelsRaw []byte `db:"labels"`
		ObservedAt time.Time `db:"observed_at"`
		CreatedAt time.Time `db:"created_at"`
	}, error)
	Ping(ctx context.Context) error
}

type Handler struct { Repo Repository }

// POST /metrics
// {
//   "service": "auth-service",
//   "name": "request_latency_ms",
//   "value": 123.4,
//   "labels": {"route":"/login","method":"POST"},
//   "observed_at": "2025-09-29T12:34:56Z"
// }
func (h Handler) Ingest(c *gin.Context) {
	var payload map[string]any
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"}); return
	}
	// default observed_at = now
	if _, ok := payload["observed_at"]; !ok { payload["observed_at"] = time.Now().UTC() }
	if err := h.Repo.InsertMetric(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}

// GET /metrics?service=auth-service&name=request_latency_ms&from=...&to=...&limit=500
func (h Handler) Query(c *gin.Context) {
	service := c.Query("service")
	name := c.Query("name")
	fromStr := c.Query("from")
	toStr := c.Query("to")
	limit := 500

	from, to := time.Now().Add(-1*time.Hour), time.Now()
	if t, err := time.Parse(time.RFC3339, fromStr); err == nil { from = t }
	if t, err := time.Parse(time.RFC3339, toStr); err == nil { to = t }

	rows, err := h.Repo.QueryMetrics(c.Request.Context(), service, name, from, to, limit)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, rows)
}

// GET /healthz
func (h Handler) Healthz(c *gin.Context) {
	if err := h.Repo.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status":"db_error", "error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"status":"ok"})
}