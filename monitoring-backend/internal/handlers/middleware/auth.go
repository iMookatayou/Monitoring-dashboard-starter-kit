package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIKEYAuth struct { Key string }

func (a APIKeyAuth) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-KEY")
		if key == "" || key != a.Key {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}