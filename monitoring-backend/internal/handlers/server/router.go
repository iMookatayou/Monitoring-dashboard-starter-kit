package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iMookatayou/monitoring-dashboard-starter-kit/monitoring-backend/internal/handlers"
	"github.com/iMookatayou/monitoring-dashboard-starter-kit/monitoring-backend/internal/middleware"
)

type Deps struct {
	Handler handlers.Handler
	ApiKey  string
}

func NewRouter(d Deps) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", d.Handler.Healthz)

	// Protected routes
	auth := middleware.APIKeyAuth{Key: d.ApiKey}
	grp := r.Group("/")
	grp.Use(auth.Handler())
	{
		grp.POST("/metrics", d.Handler.Ingest)
		grp.GET("/metrics", d.Handler.Query)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
	})

	return r
}
