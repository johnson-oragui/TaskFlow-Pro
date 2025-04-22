package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func PathMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if path == "/" || path == "/docs" || path == "/favico" || strings.HasPrefix(path, "/api/v1") {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message":     "Forbidden Path",
			"status_code": 403,
			"data":        map[string]string{},
		})
	}
}
