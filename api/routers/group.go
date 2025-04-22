package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/johnson-oragui/TaskFlow-Pro/api/routers/v1"
)

func RouterV1(r *gin.Engine) {
	v1 := r.Group("/api/v1")

	routers.RegisterAuthRoutes(v1)
	routers.RegisterUserRoutes(v1)
}
