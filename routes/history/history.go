package history

import (
	"healcationBackend/controllers/history"
	"healcationBackend/middleware"

	"github.com/gin-gonic/gin"
)

func HistoryRoutes(r *gin.Engine) {
	historyGroup := r.Group("/history", middleware.Validate())
	{
		historyGroup.GET("/", history.GetHistories)
		historyGroup.GET("/:id", history.GetHistoryDetail)
		historyGroup.DELETE("/:id", history.DeleteHistory)
	}
}
