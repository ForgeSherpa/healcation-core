package api

import (
	"healcationBackend/controllers/api"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(r *gin.Engine) {
	r.GET("/api", api.TestGemini)
}
