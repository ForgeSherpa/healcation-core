package auth

import (
	"healcationBackend/controllers/auth"
	"healcationBackend/middleware"

	"github.com/gin-gonic/gin"
)

func RoutesAuth(r *gin.Engine) {
	r.POST("/api/register", auth.Register)
	r.POST("/api/login", auth.Login)
	protected := r.Group("/")
	protected.Use(middleware.Validate())
	{
		protected.GET("/api/validate", auth.Validate)
	}
	r.POST("/api/refresh-token", auth.RefreshToken)
}
