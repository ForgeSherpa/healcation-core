package history

import (
	"healcationBackend/database"
	"healcationBackend/models"
	"math"
	"net/http"
	"strconv"
	"time"

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

type HistoryResponse struct {
	ID        uint   `json:"id"`
	Country   string `json:"country"`
	Town      string `json:"town"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Image     string `json:"image"`
}

func GetHistories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	var histories []models.History
	var totalRecords int64

	if err := database.DB.Model(&models.History{}).Count(&totalRecords).Error; err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to count histories: "+err.Error())
		return
	}

	if err := database.DB.Limit(limit).Offset(offset).Find(&histories).Error; err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to retrieve histories: "+err.Error())
		return
	}

	lastPage := int(math.Ceil(float64(totalRecords) / float64(limit)))

	if len(histories) == 0 {
		sendResponse(c, http.StatusOK, gin.H{
			"histories": []HistoryResponse{},
			"meta": gin.H{
				"current_page": page,
				"last_page":    lastPage,
			},
		}, "No histories found")
		return
	}

	responseHistories := make([]HistoryResponse, 0)
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

	sendResponse(c, http.StatusOK, gin.H{
		"histories": responseHistories,
		"meta": gin.H{
			"current_page": page,
			"last_page":    lastPage,
		},
	}, "Histories retrieved successfully")
}

func GetHistoryDetail(c *gin.Context) {
	id := c.Param("id")
	var history models.History

	if err := database.DB.First(&history, "id = ?", id).Error; err != nil {
		sendResponse(c, http.StatusNotFound, nil, "History not found")
		return
	}

	response := HistoryResponse{
		ID:        history.ID,
		Country:   history.Country,
		Town:      history.Town,
		StartDate: history.StartDate.Format(time.RFC3339),
		EndDate:   history.EndDate.Format(time.RFC3339),
		Image:     history.Image,
	}

	sendResponse(c, http.StatusOK, response, "History retrieved successfully")

}

func DeleteHistory(c *gin.Context) {
	id := c.Param("id")
	var history models.History

	if err := database.DB.First(&history, "id = ?", id).Error; err != nil {
		sendResponse(c, http.StatusNotFound, nil, "History not found")
		return
	}

	if err := database.DB.Delete(&history).Error; err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to delete history: "+err.Error())
		return
	}

	sendResponse(c, http.StatusOK, nil, "History deleted successfully")
}
