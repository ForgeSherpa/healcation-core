package planner

import (
	"healcationBackend/controllers/planner"

	"github.com/gin-gonic/gin"
)

func PlannerRoutes(r *gin.Engine) {
	plannerGroup := r.Group("/planner")
	{
		plannerGroup.GET("/popular", planner.GetPopularDestinations)
		plannerGroup.GET("/search", planner.SearchPlanner)
	}
}
