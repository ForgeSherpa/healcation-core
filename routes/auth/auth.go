package auth

import (
	"healcationBackend/controllers/auth"
	"healcationBackend/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.POST("/api/register", auth.Register)
	r.POST("/api/login", auth.Login)
	r.GET("/api/validate", middleware.Validate(), auth.Validate)
	r.POST("/api/refresh-token", auth.RefreshToken)
}
