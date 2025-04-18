package routers

import "github.com/gin-gonic/gin"

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login")
	}
}
