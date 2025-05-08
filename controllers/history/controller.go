package history

import (
	"encoding/json"
	"healcationBackend/database"
	"healcationBackend/models"
	"math"
	"net/http"
	"strconv"
	"strings"
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
	uidValue, exists := c.Get("userID")
	if !exists {
		sendResponse(c, http.StatusUnauthorized, nil, "Unauthorized: user not found in context")
		return
	}
	userID, ok := uidValue.(uint)
	if !ok {
		sendResponse(c, http.StatusInternalServerError, nil, "Invalid user ID type")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	search := strings.TrimSpace(c.Query("search"))

	var pattern string
	if search != "" {
		pattern = "%" + strings.ToLower(search) + "%"
	}

	var totalRecords int64
	countQ := database.DB.Model(&models.History{}).
		Where("user_id = ?", userID)
	if search != "" {
		countQ = countQ.Where("LOWER(town) LIKE ?", pattern)
	}
	if err := countQ.Count(&totalRecords).Error; err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to count histories: "+err.Error())
		return
	}

	var histories []models.History
	dataQ := database.DB.
		Where("user_id = ?", userID)
	if search != "" {
		dataQ = dataQ.Where("LOWER(town) LIKE ?", pattern)
	}
	if err := dataQ.
		Limit(limit).
		Offset(offset).
		Find(&histories).Error; err != nil {
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
		firstImage := ""
		var imgs []string
		if err := json.Unmarshal([]byte(h.Image), &imgs); err == nil {
			if len(imgs) > 0 {
				firstImage = imgs[0]
			}
		}

		responseHistories = append(responseHistories, HistoryResponse{
			ID:        h.ID,
			Country:   h.Country,
			Town:      h.Town,
			StartDate: h.StartDate.Format(time.RFC3339),
			EndDate:   h.EndDate.Format(time.RFC3339),
			Image:     firstImage,
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
	uidValue, exists := c.Get("userID")
	if !exists {
		sendResponse(c, http.StatusUnauthorized, nil, "Unauthorized: user not found in context")
		return
	}
	userID, ok := uidValue.(uint)
	if !ok {
		sendResponse(c, http.StatusInternalServerError, nil, "Invalid user ID type")
		return
	}

	id := c.Param("id")
	var history models.History

	if err := database.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&history).Error; err != nil {
		sendResponse(c, http.StatusNotFound, nil, "History not found or access denied")
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
	uidValue, exists := c.Get("userID")
	if !exists {
		sendResponse(c, http.StatusUnauthorized, nil, "Unauthorized: user not found in context")
		return
	}
	userID, ok := uidValue.(uint)
	if !ok {
		sendResponse(c, http.StatusInternalServerError, nil, "Invalid user ID type")
		return
	}

	id := c.Param("id")
	var history models.History

	if err := database.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&history).Error; err != nil {
		sendResponse(c, http.StatusNotFound, nil, "History not found or access denied")
		return
	}

	if err := database.DB.Delete(&history).Error; err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to delete history: "+err.Error())
		return
	}
	sendResponse(c, http.StatusOK, nil, "History deleted successfully")
}
