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

	// public
	mux.HandleFunc("/healthz", d.Handler.Healthz)

	// protected
	metricsIngest := http.HandlerFunc(d.Handler.Ingest)
	metricsQuery := http.HandlerFunc(d.Handler.Query)

	protected := func(h http.Handler) http.Handler {
		return middleware.Chain(h, middleware.APIKey(d.APIKey))
	}

	mux.Handle("/metrics", protected(metricsQuery))   // GET
	mux.Handle("/metrics/", protected(metricsIngest)) // POST (method checked client-side)

	return mux
}
