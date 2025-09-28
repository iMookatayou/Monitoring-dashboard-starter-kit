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
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSL)
	repo, err := repository.New(dsn)
	if err != nil {
		log.Fatalf("db connect/migrate: %v", err)
	}

	h := handlers.Handler{Repo: repo}
	r := server.NewRouter(server.Deps{Handler: h, ApiKey: cfg.ApiKey})

	srv := &http.Server{Addr: ":" + cfg.AppPort, Handler: r, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second}
	log.Printf("listening on :%s", cfg.AppPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
