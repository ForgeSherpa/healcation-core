package api

import (
	"healcationBackend/database"
	"healcationBackend/routes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	database.LoadEnvVariables()
	database.Connect()

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard // Hindari conflict log Vercel

	router := gin.Default()
	routes.Routes(router)

	router.ServeHTTP(w, r)
}
