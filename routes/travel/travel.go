package travel

import (
	"healcationBackend/controllers/travel"
	"healcationBackend/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	travelGroup := r.Group("/travel", middleware.Validate())
	{
		travelGroup.POST("/places", travel.GetPlaces)
		travelGroup.POST("/places-detail/:name", travel.GetPlaceDetail)
		travelGroup.POST("/select-place", travel.SelectPlace)
	}
	r.POST("/timeline", middleware.Validate(), travel.Timeline)
}
