package routes

import (
	"healcationBackend/routes/auth"
	"healcationBackend/routes/history"
	"healcationBackend/routes/planner"
	"healcationBackend/routes/profile"
	"healcationBackend/routes/travel"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	auth.RoutesAuth(r)
	history.HistoryRoutes(r)
	profile.ProfileRoutes(r)
	planner.PlannerRoutes(r)
	travel.TravelRoutes(r)
}
