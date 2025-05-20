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

type DayPlace struct {
	Date string         `json:"date"`
	Data []placeVisited `json:"data"`
}

type placeVisited struct {
	Image    []string `json:"image"`
	Landmark string   `json:"landmark"`
	RoadName string   `json:"roadName"`
	Time     string   `json:"time"`
	Town     string   `json:"town"`
	Type     string   `json:"type"`
}

type SelectPlaceRequest struct {
	Country      string     `json:"country"`
	Town         string     `json:"town"`
	StartDate    string     `json:"startDate"`
	EndDate      string     `json:"endDate"`
	PlaceVisited []DayPlace `json:"placeVisited"`
	Budget       string     `json:"budget"`
}

type PlaceData struct {
	Country      string     `json:"country"`
	Town         string     `json:"town"`
	StartDate    string     `json:"startDate"`
	EndDate      string     `json:"endDate"`
	Budget       string     `json:"budget"`
	PlaceVisited []DayPlace `json:"placeVisited"`
}

type placeVisitedSimple struct {
	Image    string `json:"image"`
	Landmark string `json:"landmark"`
	RoadName string `json:"roadName"`
	Time     string `json:"time"`
	Town     string `json:"town"`
	Type     string `json:"type"`
}

type SelectPlaceResponse struct {
	Message string                          `json:"message"`
	Data    map[string][]placeVisitedSimple `json:"data"`
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

	var req SelectPlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendResponse(c, http.StatusBadRequest, nil, "Invalid request format: "+err.Error())
		return
	}

	visitMap := make(map[string][]placeVisitedSimple)
	for _, day := range req.PlaceVisited {
		if day.Date == "" || len(day.Data) == 0 {
			continue
		}
		for _, v := range day.Data {
			// pick first image or empty
			img := ""
			if len(v.Image) > 0 {
				img = v.Image[0]
			}
			simple := placeVisitedSimple{
				Image:    img,
				Landmark: v.Landmark,
				RoadName: v.RoadName,
				Time:     v.Time,
				Town:     v.Town,
				Type:     v.Type,
			}
			visitMap[day.Date] = append(visitMap[day.Date], simple)
		}
	}

	firstImage := ""
	for _, visits := range visitMap {
		if len(visits) > 0 {
			firstImage = visits[0].Image
			break
		}
	}

	tlJSON, err := json.Marshal(visitMap)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to marshal placeVisited: "+err.Error())
		return
	}

	history := models.History{
		UserID:    userID,
		Country:   req.Country,
		Town:      req.Town,
		StartDate: parseDate(req.StartDate),
		EndDate:   parseDate(req.EndDate),
		Budget:    req.Budget,
		Timelines: string(tlJSON),
		Image:     firstImage,
	}

	if err := database.DB.Create(&history).Error; err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to save history: "+err.Error())
		return
	}

	resp := SelectPlaceResponse{
		Message: "Done! Enjoy your vacation!",
		Data:    visitMap,
	}

	sendResponse(c, http.StatusOK, resp, "Place selection saved successfully")
}
