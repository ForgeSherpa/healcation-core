package travel

import (
	"healcationBackend/controllers/travel"

	"github.com/gin-gonic/gin"
)

func TravelRoutes(r *gin.Engine) {
	travelGroup := r.Group("/travel")
	{
		travelGroup.GET("/", travel.GetPreferences)
		travelGroup.GET("/select_place", travel.GetSelectPlaces)
	}

}
