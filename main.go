package main

import (
	"healcationBackend/database"
	"healcationBackend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.LoadEnvVariables()
	database.Connect()

	r := gin.Default()
	routes.Routes(r)

	r.Run(":3000")
}
