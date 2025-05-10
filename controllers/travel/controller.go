package travel

import (
	"encoding/json"
	"errors"
	"fmt"
	"healcationBackend/database"
	"healcationBackend/models"
	"healcationBackend/pkg/services"
	"net/http"
	"net/url"
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

func GetPlaceDetail(c *gin.Context) {
	placeName, err := url.QueryUnescape(c.Param("name"))
	if err != nil {
		sendResponse(c, http.StatusBadRequest, nil, "Invalid place name")
		return
	}

	var request struct {
		Country string `json:"country"`
		City    string `json:"city"`
		Type    string `json:"type"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		sendResponse(c, http.StatusBadRequest, nil, "Invalid request format: "+err.Error())
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

	placeDetail, err := aiSvc.GetPlaceDetail(placeName, request.Type, request.Country, request.City)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to fetch place details from AI Service: "+err.Error())
		return
	}

	sendResponse(c, http.StatusOK, gin.H{"place_detail": placeDetail}, "Place detail retrieved successfully")
}

type TimelineRequest struct {
	Accomodation string `json:"accomodation"`
	Town         string `json:"town"`
	Country      string `json:"country"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	Places       []struct {
		Name      string `json:"name"`
		TimeOfDay string `json:"timeOfDay"`
	} `json:"places"`
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
		request.Places,
	)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Gagal mendapatkan response dari AI Service: "+err.Error())
		return
	}

	sendResponse(c, http.StatusOK, gin.H{"timeline": response}, "Timeline retrieved successfully")
}

type SelectPlaceRequest struct {
	Country      string                      `json:"country"`
	Town         string                      `json:"town"`
	StartDate    string                      `json:"startDate"`
	EndDate      string                      `json:"endDate"`
	Accomodation string                      `json:"accomodation"`
	Title        string                      `json:"title"`
	Timelines    map[string][]TimelineDetail `json:"timelines"`
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

	userIDValue, exists := c.Get("userID")
	if !exists {
		sendResponse(c, http.StatusUnauthorized, nil, "Unauthorized: user not found in context")
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		sendResponse(c, http.StatusInternalServerError, nil, "Invalid user ID type")
		return
	}

	var request SelectPlaceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		sendResponse(c, http.StatusBadRequest, nil, "Invalid request format: "+err.Error())
		return
	}

	var allImages []string
	for _, dayDetails := range request.Timelines {
		for _, detail := range dayDetails {
			allImages = append(allImages, detail.Image)
			fmt.Println("All Images:", allImages)
		}
	}
	imageJSON, err := json.Marshal(allImages)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to process images")
		return
	}

	history := models.History{
		UserID:    userID,
		Country:   request.Country,
		Town:      request.Town,
		StartDate: parseDate(request.StartDate),
		EndDate:   parseDate(request.EndDate),
		Image:     string(imageJSON),
	}

	if err := database.DB.Create(&history).Error; err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to save history")
		return
	}

	response := SelectPlaceResponse{
		Message: "Done! Enjoy your vacation!",
		Data: PlaceData{
			Country:   request.Country,
			Town:      request.Town,
			Title:     request.Title,
			StartDate: request.StartDate,
			EndDate:   request.EndDate,
			SelectedAccomodation: []AccomodationDetail{
				{
					Name:     request.Accomodation,
					RoadName: "",
				},
			},
			Timeline: request.Timelines,
		},
	}

	sendResponse(c, http.StatusOK, response, "Place selection saved successfully")
}
