package planner

import (
	"healcationBackend/controllers/planner"
	"healcationBackend/middleware"

	"github.com/gin-gonic/gin"
)

func PlannerRoutes(r *gin.Engine) {
	publicPlanner := r.Group("/planner")
	{
		publicPlanner.GET("/popular", planner.GetPopularDestinations)
	}

	privatePlanner := r.Group("/planner")
	privatePlanner.Use(middleware.Validate())
	{
		privatePlanner.GET("/search", planner.SearchPlanner)
	}
}
