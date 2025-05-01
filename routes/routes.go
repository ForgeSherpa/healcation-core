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

	auth.Routes(r)
	history.Routes(r)
	profile.Routes(r)
	planner.Routes(r)
	travel.Routes(r)
}
