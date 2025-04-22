package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DefaultRouters(r *gin.Engine) {
	r.GET("", func(c *gin.Context) { c.JSON(200, gin.H{"message": "HOME"}) })
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204) // or serve actual favicon
	})
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
