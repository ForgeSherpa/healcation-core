package history

import (
	"healcationBackend/controllers/history"
	"healcationBackend/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	historyGroup := r.Group("/history", middleware.Validate())
	{
		historyGroup.GET("/", history.GetHistories)
		historyGroup.GET("/:id", history.GetHistoryDetail)
		historyGroup.GET("/:id/place", history.GetHistoryDetailPlace)
		historyGroup.DELETE("/:id", history.DeleteHistory)
	}
}
