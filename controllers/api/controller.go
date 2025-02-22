package api

import (
	"healcationBackend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestGemini(c *gin.Context) {
	city := c.Query("city")
	geminiData, err := services.FetchFromGeminiAPI(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, geminiData)
}
