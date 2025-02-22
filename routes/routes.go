package routes

import (
	"healcationBackend/routes/api"
	"healcationBackend/routes/auth"
	"healcationBackend/routes/home"
	"healcationBackend/routes/planner"
	"healcationBackend/routes/timeline"
	"healcationBackend/routes/travel"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	auth.RoutesAuth(r)
	home.HomeRoutes(r)
	planner.PlannerRoutes(r)
	travel.TravelRoutes(r)
	timeline.TimelineRoutes(r)
	api.ApiRoutes(r)
}
