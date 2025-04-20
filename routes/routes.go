package routes

import (
	"healcationBackend/routes/auth"
	"healcationBackend/routes/history"
	"healcationBackend/routes/planner"
	"healcationBackend/routes/profile"
	"healcationBackend/routes/travel"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API Healcation jalan!"})
	})

	auth.RoutesAuth(r)
	history.HistoryRoutes(r)
	profile.ProfileRoutes(r)
	planner.PlannerRoutes(r)
	travel.TravelRoutes(r)
}
