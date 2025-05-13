package planner

import (
	"healcationBackend/controllers/planner"
	"healcationBackend/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	privatePlanner := r.Group("/planner")
	privatePlanner.Use(middleware.Validate())
	{
		privatePlanner.GET("/search", planner.SearchPlanner)
		privatePlanner.GET("/popular", planner.SearchPlanner)
	}
}
