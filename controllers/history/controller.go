package history

import (
	"healcationBackend/database"
	"healcationBackend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HistoryResponse struct {
	ID        uint   `json:"id"`
	Country   string `json:"country"`
	Town      string `json:"town"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Image     string `json:"image"`
}

func GetHistories(c *gin.Context) {
	var histories []models.History
	database.DB.Find(&histories)

	var responseHistories []HistoryResponse
	for _, h := range histories {
		responseHistories = append(responseHistories, HistoryResponse{
			ID:        h.ID,
			Country:   h.Country,
			Town:      h.Town,
			StartDate: h.StartDate.Format(time.RFC3339),
			EndDate:   h.EndDate.Format(time.RFC3339),
			Image:     h.Image,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"histories": responseHistories,
		"meta": gin.H{
			"current_page": 1,
			"last_page":    1,
		},
	})
}

func GetHistoryDetail(c *gin.Context) {
	id := c.Param("id")
	var history models.History

	if err := database.DB.Preload("SelectedAccomodation").Preload("SelectedPlaces").First(&history, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}

func DeleteHistory(c *gin.Context) {
	id := c.Param("id")
	var history models.History

	if err := database.DB.First(&history, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	database.DB.Where("history_id = ?", id).Delete(&models.SelectedAccomodation{})
	database.DB.Where("history_id = ?", id).Delete(&models.SelectedPlace{})

	if err := database.DB.Delete(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "History deleted!"})
}
