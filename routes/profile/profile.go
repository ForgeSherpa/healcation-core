package profile

import (
	"healcationBackend/controllers/profile"
	"healcationBackend/middleware"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(r *gin.Engine) {
	profileGroup := r.Group("/api/profile", middleware.Validate())
	{
		profileGroup.GET("/", profile.GetProfile)
		profileGroup.PUT("/", profile.UpdateProfile)
	}
}
