package home

import (
	"github.com/gin-gonic/gin"
	"healcationBackend/controllers/home"
	"healcationBackend/middleware"
)

func HomeRoutes(r *gin.Engine) {
	homeGroup := r.Group("/home")
	{
		homeGroup.GET("/history", middleware.Validate(), home.GetHistory)
		homeGroup.GET("/popular", home.GetPopularDestinations)
	}
}
