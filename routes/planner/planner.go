package planner

import (
	"github.com/gin-gonic/gin"
	"healcationBackend/controllers/planner"
)

func PlannerRoutes(r *gin.Engine) {
	plannerGroup := r.Group("/planner")
	{
		plannerGroup.GET("/search", planner.SearchPlaces)
		plannerGroup.GET("/popular", planner.GetPopularPlaces)
	}

}
