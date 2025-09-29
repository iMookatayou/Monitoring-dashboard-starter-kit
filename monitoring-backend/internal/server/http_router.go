package server

import (
	"net/http"

	"github.com/iMookatayou/monitoring-dashboard-starter-kit/monitoring-backend/internal/handlers"
	"github.com/iMookatayou/monitoring-dashboard-starter-kit/monitoring-backend/internal/middleware"
)

type Deps struct {
	Handler handlers.HTTP
	APIKey  string
}

func NewMux(d Deps) *http.ServeMux {
	mux := http.NewServeMux()

	// --- public ---
	mux.HandleFunc("/healthz", d.Handler.Healthz)

	// --- protected helper ---
	protected := func(h http.Handler) http.Handler {
		return middleware.Chain(h, middleware.APIKey(d.APIKey))
	}

	// --- /metrics: แยกตาม method ที่ path เดียว ---
	mux.Handle("/metrics", protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			d.Handler.Query(w, r)
		case http.MethodPost:
			d.Handler.Ingest(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// (ถ้าชอบให้ /metrics/ ก็ชี้มาที่ตัวเดียวกันได้)
	mux.Handle("/metrics/", protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			d.Handler.Query(w, r)
		case http.MethodPost:
			d.Handler.Ingest(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	return mux
}
