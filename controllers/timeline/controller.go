package timeline

import (
	"github.com/gin-gonic/gin"
	"healcationBackend/database"
	"healcationBackend/models"
	"net/http"
)

func GetTimelines(c *gin.Context) {
	var timelines []models.Timeline

	// Preload PlaceVisited relationship
	if err := database.DB.Preload("PlaceVisited").Find(&timelines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve timelines"})
		return
	}

	// Format response
	response := make([]map[string]interface{}, len(timelines))
	for i, timeline := range timelines {
		placeVisited := make([]map[string]interface{}, len(timeline.PlaceVisited))

		for j, place := range timeline.PlaceVisited {
			// Append PlaceVisited data into response
			placeVisited[j] = map[string]interface{}{
				"type":     place.Type,
				"landmark": place.Landmark,
				"roadName": place.RoadName,
				"town":     place.Town,
				"time":     place.Time,
				"image":    place.Images,
			}
		}

		response[i] = map[string]interface{}{
			"id":           timeline.ID,
			"town":         timeline.Town,
			"country":      timeline.Country,
			"budget":       timeline.Budget,
			"startDate":    timeline.StartDate,
			"endDate":      timeline.EndDate,
			"placeVisited": placeVisited,
		}
	}

	c.JSON(http.StatusOK, gin.H{"timelines": response})
}
