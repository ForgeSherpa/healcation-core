package planner

import (
	"github.com/gin-gonic/gin"
	"healcationBackend/database"
	"healcationBackend/models"
	"net/http"
)

// SearchPlaces handles searching for places by town or country
func SearchPlaces(c *gin.Context) {
	// Get the search query from the request
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query cannot be empty"})
		return
	}

	// Perform the database query to search for matching places
	var results []models.Popular
	searchQuery := "%" + query + "%" // For partial matching
	if err := database.DB.Where("town LIKE ? OR country LIKE ?", searchQuery, searchQuery).Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
		return
	}

	// Format the response
	response := make([]map[string]interface{}, len(results))
	for i, place := range results {
		response[i] = map[string]interface{}{
			"id":      place.ID,
			"town":    place.Town,
			"country": place.Country,
			"image":   place.Image,
		}
	}

	c.JSON(http.StatusOK, gin.H{"results": response})
}

// GetPopularPlaces handles fetching popular destinations for the planner
func GetPopularPlaces(c *gin.Context) {
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
