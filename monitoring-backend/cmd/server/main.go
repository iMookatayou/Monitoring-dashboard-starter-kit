package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/iMookatayou/monitoring-dashboard-starter-kit/monitoring-backend/internal/handlers"
	"github.com/iMookatayou/monitoring-dashboard-starter-kit/monitoring-backend/internal/repository"
	"github.com/iMookatayou/monitoring-dashboard-starter-kit/monitoring-backend/internal/server"
	"github.com/iMookatayou/monitoring-dashboard-starter-kit/monitoring-backend/pkg/config"
)

func main() {
	cfg := config.Load()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSL,
	)

	repo, err := repository.New(dsn)
	if err != nil {
		log.Fatalf("db connect/migrate: %v", err)
	}

	h := handlers.HTTP{Repo: repo}
	mux := server.NewMux(server.Deps{Handler: h, APIKey: cfg.ApiKey})

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      middlewareChain(mux), // CORS + logging + recover
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("listening on :%s", cfg.AppPort)
	log.Fatal(srv.ListenAndServe())
}

// ---- global middlewares ----

// CORS ต้องอยู่นอกสุด เพื่อให้ preflight OPTIONS ผ่านได้ก่อน middleware อื่น
func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func middlewareChain(h http.Handler) http.Handler {
	return middlewareCORS(
		middlewareRecover(
			middlewareLogging(h),
		),
	)
}

func middlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func middlewareRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
