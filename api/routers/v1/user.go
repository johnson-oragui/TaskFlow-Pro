package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/users")
	{
		auth.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message":     "success",
				"status_code": 200,
				"data":        []string{},
			})
		})
	}
}
