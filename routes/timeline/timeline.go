package timeline

import (
	"github.com/gin-gonic/gin"
	"healcationBackend/controllers/timeline"
)

func TimelineRoutes(r *gin.Engine) {
	timelineGroup := r.Group("/timeline")
	{
		timelineGroup.GET("/", timeline.GetTimelines)
	}

}
