package planner

import (
	"errors"
	"healcationBackend/pkg/services"
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

	aiSvc, err := services.NewAIService()
	if err != nil {
		if errors.Is(err, services.ErrGeminiUnavailable) {
			services.HandleGeminiUnavailable(c.Writer, err)
		} else {
			sendResponse(c, http.StatusInternalServerError, nil, "Failed to initialize AI service: "+err.Error())
		}
		return
	}

	results, err := aiSvc.Search(query)
	if err != nil {
		sendResponse(c, http.StatusBadGateway, nil, "Failed to fetch data from AI Service: "+err.Error())
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
			"image":   "https://lp-cms-production.imgix.net/2023-01/GettyImages-827446284.jpg?w=1095&fit=crop&crop=faces%2Cedges&auto=format&q=75",
		},
		{
			"id":      "2",
			"town":    "Yogyakarta",
			"country": "Indonesia",
			"image":   "https://api2.kemenparekraf.go.id/storage/app/resources/image_artikel/Sumbu%20Kosmologis%20Yogyakarta_Shutterstock%201021926103_Creativa%20Images.jpg",
		},
		{
			"id":      "3",
			"town":    "Lombok",
			"country": "Indonesia",
			"image":   "https://img.jakpost.net/c/2016/12/15/2016_12_15_17857_1481791864._large.jpg",
		},
	}

	sendResponse(c, http.StatusOK, gin.H{"popular_destinations": popularDestinations}, "Popular destinations retrieved successfully")
}
