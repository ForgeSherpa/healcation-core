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
	userID, _ := c.Get("userID")

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
		lowered := strings.ToLower(search)
		countQ = countQ.Where("LOWER(town) LIKE ?", "%"+lowered+"%")
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

type AccomodationDetail struct {
	Name     string `json:"name"`
	RoadName string `json:"roadName"`
}

type TimelineDetail struct {
	Image    string `json:"image"`
	Landmark string `json:"landmark"`
	RoadName string `json:"roadName"`
	Time     string `json:"time"`
	Town     string `json:"town"`
	Type     string `json:"type"`
}

func GetHistoryDetail(c *gin.Context) {
	userID, _ := c.Get("userID")
	var h models.History
	idParam := c.Param("id")
	if err := database.DB.
		Where("id = ? AND user_id = ?", idParam, userID).
		First(&h).Error; err != nil {
		sendResponse(c, http.StatusNotFound, nil, "Not found")
		return
	}

	var accoms []AccomodationDetail
	json.Unmarshal([]byte(h.Accommodations), &accoms)

	var timelines map[string][]TimelineDetail
	json.Unmarshal([]byte(h.Timelines), &timelines)

	type PlaceData struct {
		Type     string `json:"type"`
		Landmark string `json:"landmark"`
		RoadName string `json:"roadName"`
		Town     string `json:"town"`
		Time     string `json:"time"`
		Image    string `json:"image"`
	}
	type DateGroup struct {
		Date string      `json:"date"`
		Data []PlaceData `json:"data"`
	}

	var placeVisited []DateGroup
	for date, items := range timelines {
		var list []PlaceData
		for _, it := range items {
			list = append(list, PlaceData{
				Type:     it.Type,
				Landmark: it.Landmark,
				RoadName: it.RoadName,
				Town:     it.Town,
				Time:     it.Time,
				Image:    it.Image,
			})
		}
		placeVisited = append(placeVisited, DateGroup{Date: date, Data: list})
	}

	resp := struct {
		ID           uint        `json:"id"`
		Budget       string      `json:"budget"`
		Town         string      `json:"town"`
		Country      string      `json:"country"`
		StartDate    string      `json:"startDate"`
		EndDate      string      `json:"endDate"`
		PlaceVisited []DateGroup `json:"placeVisited"`
	}{
		ID:           h.ID,
		Budget:       h.Budget,
		Town:         h.Town,
		Country:      h.Country,
		StartDate:    h.StartDate.Format(time.RFC3339Nano),
		EndDate:      h.EndDate.Format(time.RFC3339Nano),
		PlaceVisited: placeVisited,
	}

	sendResponse(c, http.StatusOK, resp, "History detail retrieved")
}

func DeleteHistory(c *gin.Context) {
	userID, _ := c.Get("userID")

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
