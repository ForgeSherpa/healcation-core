package main

import (
	"healcationBackend/database"
	"healcationBackend/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	database.LoadEnvVariables()
	database.Connect()
}
func main() {
	r := gin.Default()
	routes.Routes(r)
	r.Run()

	// r.Run("0.0.0.0:8080")
}
