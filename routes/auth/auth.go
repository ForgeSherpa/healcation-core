package auth

import (
	"github.com/gin-gonic/gin"
	"healcationBackend/controllers/auth"
	"healcationBackend/middleware"
)

func RoutesAuth(r *gin.Engine) {
	r.POST("/api/register", auth.Register)
	r.POST("/api/login", auth.Login)
	protected := r.Group("/")
	protected.Use(middleware.Validate())
	{
		protected.GET("/api/validate", auth.Validate)
	}
	r.POST("/api/logout", auth.Logout)
	r.POST("/api/refresh-token", auth.RefreshToken)
}
