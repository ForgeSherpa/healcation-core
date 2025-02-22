package main

import (
	"github.com/gin-gonic/gin"
	"healcationBackend/database"
	"healcationBackend/routes"
)

func init() {
	database.LoadEnvVariables()
	database.Connect()
}
func main() {
	r := gin.Default()
	routes.Routes(r)
	r.Run()
}
