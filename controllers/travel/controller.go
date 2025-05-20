package travel

import (
	"encoding/json"
	"errors"
	"healcationBackend/database"
	"healcationBackend/models"
	"healcationBackend/pkg/services"
	"net/http"
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

func GetPlaces(c *gin.Context) {
	var request struct {
		Preferences []string `json:"preferences"`
		Country     string   `json:"country"`
		Town        string   `json:"town"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		sendResponse(c, http.StatusBadRequest, nil, "Format JSON salah: "+err.Error())
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

	placesData, err := aiSvc.GetPlaces(request.Preferences, request.Country, request.Town)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Gagal mengambil data dari AI Service: "+err.Error())
		return
	}

	sendResponse(c, http.StatusOK, gin.H{"places": placesData}, "Places retrieved successfully")
}

type LandmarkDetailResponse struct {
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

func GetPlaceDetail(c *gin.Context) {
	var req struct {
		Type     string `json:"type" binding:"required"`
		Landmark string `json:"landmark" binding:"required"`
		Town     string `json:"town" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		sendResponse(c, http.StatusBadRequest, nil,
			"Invalid or missing JSON fields: "+err.Error())
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

	detailMap, err := aiSvc.GetPlaceDetail(req.Type, req.Landmark, req.Town)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to fetch landmark details from AI Service: "+err.Error())
		return
	}

	resp := LandmarkDetailResponse{
		Description: detailMap["description"].(string),
		Images:      detailMap["images"].([]string),
	}

	sendResponse(c, http.StatusOK, resp, "Place detail retrieved successfully")
}

type TimelineRequest struct {
	Accomodation  string                   `json:"accomodation"`
	Town          string                   `json:"town"`
	Country       string                   `json:"country"`
	StartDate     string                   `json:"startDate"`
	EndDate       string                   `json:"endDate"`
	SelectedPlace []services.SelectedPlace `json:"selectedPlace"`
}

func Timeline(c *gin.Context) {
	var request TimelineRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		sendResponse(c, http.StatusBadRequest, nil, "Format JSON salah: "+err.Error())
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

	response, err := aiSvc.GetTimeline(
		request.Accomodation,
		request.Town,
		request.Country,
		request.StartDate,
		request.EndDate,
		request.SelectedPlace,
	)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Gagal mendapatkan response dari AI Service: "+err.Error())
		return
	}

	sendResponse(c, http.StatusOK, response, "Timeline retrieved successfully")
}

type SelectPlaceRequest struct {
	Country        string                      `json:"country"`
	Town           string                      `json:"town"`
	StartDate      string                      `json:"startDate"`
	EndDate        string                      `json:"endDate"`
	Accommodations []AccomodationDetail        `json:"accommodations"`
	Title          string                      `json:"title"`
	Timelines      map[string][]TimelineDetail `json:"timelines"`
	Budget         string                      `json:"budget"`
}

type TimelineDetail struct {
	Image    string `json:"image"`
	Landmark string `json:"landmark"`
	RoadName string `json:"roadName"`
	Time     string `json:"time"`
	Town     string `json:"town"`
	Type     string `json:"type"`
}

type SelectPlaceResponse struct {
	Message string    `json:"message"`
	Data    PlaceData `json:"data"`
}

type PlaceData struct {
	Country              string                      `json:"country"`
	Town                 string                      `json:"town"`
	Title                string                      `json:"title"`
	StartDate            string                      `json:"startDate"`
	EndDate              string                      `json:"endDate"`
	Budget               string                      `json:"budget"`
	SelectedAccomodation []AccomodationDetail        `json:"selectedAccomodation"`
	Timeline             map[string][]TimelineDetail `json:"timeline"`
}

type AccomodationDetail struct {
	Name     string `json:"name"`
	RoadName string `json:"roadName"`
}

func parseDate(dateStr string) time.Time {
	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now()
	}
	return parsedTime
}

func SelectPlace(c *gin.Context) {
	userIDValue, _ := c.Get("userID")
	userID, _ := userIDValue.(uint)

	var request SelectPlaceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		sendResponse(c, http.StatusBadRequest, nil, "Invalid request format: "+err.Error())
		return
	}

	var firstImage string
	if details, ok := request.Timelines["1"]; ok && len(details) > 0 {
		firstImage = details[0].Image
	} else {
		for _, details := range request.Timelines {
			if len(details) > 0 {
				firstImage = details[0].Image
				break
			}
		}
	}
	accJSON, err := json.Marshal(request.Accommodations)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to marshal accommodations")
		return
	}

	tlJSON, err := json.Marshal(request.Timelines)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to marshal timelines")
		return
	}

	history := models.History{
		UserID:         userID,
		Country:        request.Country,
		Town:           request.Town,
		Title:          request.Title,
		StartDate:      parseDate(request.StartDate),
		EndDate:        parseDate(request.EndDate),
		Budget:         request.Budget,
		Accommodations: string(accJSON),
		Timelines:      string(tlJSON),
		Image:          firstImage,
	}

	if err := database.DB.Create(&history).Error; err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to save history")
		return
	}

	response := SelectPlaceResponse{
		Message: "Done! Enjoy your vacation!",
		Data: PlaceData{
			Country:              request.Country,
			Town:                 request.Town,
			Title:                request.Title,
			StartDate:            request.StartDate,
			EndDate:              request.EndDate,
			Budget:               request.Budget,
			SelectedAccomodation: request.Accommodations,
			Timeline:             request.Timelines,
		},
	}

	sendResponse(c, http.StatusOK, response, "Place selection saved successfully")
}
