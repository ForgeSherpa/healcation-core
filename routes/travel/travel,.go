package travel

import (
	"healcationBackend/controllers/travel"

	"github.com/gin-gonic/gin"
)

func TravelRoutes(r *gin.Engine) {
	travelGroup := r.Group("/travel")
	{
		travelGroup.POST("/places", travel.GetPlaces)
		travelGroup.POST("/places-detail/:name", travel.GetPlaceDetail)
		travelGroup.POST("/select-place", travel.SelectPlace)
	}
	r.POST("/timeline", travel.Timeline)
}
