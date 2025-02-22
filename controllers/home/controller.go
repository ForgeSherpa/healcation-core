package home

import (
	"github.com/gin-gonic/gin"
	"healcationBackend/database"
	"healcationBackend/models"
	"net/http"
)

func GetHistory(c *gin.Context) {
	// Retrieve the userID from the context, which was set by the Validate middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch the user's history data from the database based on the userID
	var histories []models.History
	if err := database.DB.Where("user_id = ?", userID).Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history data"})
		return
	}

	// Format the response into a consistent JSON format
	response := make([]map[string]interface{}, len(histories))
	for i, history := range histories {
		response[i] = map[string]interface{}{
			"id":        history.ID,
			"town":      history.Town,
			"country":   history.Country,
			"startDate": history.StartDate,
			"endDate":   history.EndDate,
			"image":     history.Image,
		}
	}

	// Return the history data in the response
	c.JSON(http.StatusOK, gin.H{"histories": response})
}

func GetPopularDestinations(c *gin.Context) {
	var popularDestinations []models.Popular
	if err := database.DB.Find(&popularDestinations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve popular destinations"})
		return
	}

	response := make([]map[string]interface{}, len(popularDestinations))
	for i, popular := range popularDestinations {
		response[i] = map[string]interface{}{
			"id":      popular.ID,
			"town":    popular.Town,
			"country": popular.Country,
			"image":   popular.Image,
		}
	}

	c.JSON(http.StatusOK, gin.H{"popular_destinations": response})
}
