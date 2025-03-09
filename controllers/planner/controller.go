package planner

import (
	"healcationBackend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

func sendResponse(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, Response{
		Status:  status,
		Data:    data,
		Message: message,
	})
}

func SearchPlanner(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		sendResponse(c, http.StatusBadRequest, nil, "Query parameter is required")
		return
	}

	results, err := services.SearchGemini(query)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to fetch data from Gemini AI: "+err.Error())
		return
	}

	sendResponse(c, http.StatusOK, gin.H{"results": results}, "Search results retrieved successfully")
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

	sendResponse(c, http.StatusOK, gin.H{"popular_destinations": popularDestinations}, "Popular destinations retrieved successfully")
}
