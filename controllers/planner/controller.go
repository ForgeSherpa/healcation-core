package planner

import (
	"healcationBackend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchPlanner(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	results, err := services.SearchGemini(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from Gemini AI"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
	})
}

func GetPopularDestinations(c *gin.Context) {
	popularDestinations := []map[string]string{
		{
			"id":      "1",
			"town":    "Bali",
			"country": "Indonesia",
			"image":   "https://example.com/bali.jpg",
		},
		{
			"id":      "2",
			"town":    "Yogyakarta",
			"country": "Indonesia",
			"image":   "https://example.com/jogja.jpg",
		},
		{
			"id":      "3",
			"town":    "Lombok",
			"country": "Indonesia",
			"image":   "https://example.com/lombok.jpg",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"popular_destinations": popularDestinations,
	})
}
