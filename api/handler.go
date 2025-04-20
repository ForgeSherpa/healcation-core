package api

import (
	"healcationBackend/database"
	"healcationBackend/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler dipanggil Vercel setiap request
func Handler(w http.ResponseWriter, r *http.Request) {
	// Muat variabel env & koneksi DB sekali per cold start
	database.LoadEnvVariables()
	database.Connect()

	// Buat router Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routes.Routes(router)

	// Serve HTTP dengan Gin
	router.ServeHTTP(w, r)
}
