package travel

import (
	"fmt"
	"healcationBackend/database"
	"healcationBackend/models"
	"healcationBackend/services"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

func GetPlaces(c *gin.Context) {
	var request struct {
		Preferences []string `json:"preferences"`
		Country     string   `json:"country"`
		Town        string   `json:"town"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println("❌ Format JSON salah:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON salah"})
		return
	}

	fmt.Println("✅ Request dari Postman:", request)

	placesData, err := services.GetPlacesFromGemini(request.Preferences, request.Country, request.Town)
	if err != nil {
		fmt.Println("❌ Gagal mengambil data dari Gemini:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dari Gemini API"})
		return
	}

	c.JSON(http.StatusOK, placesData)
}

func GetPlaceDetail(c *gin.Context) {
	placeName, err := url.QueryUnescape(c.Param("name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place name"})
		return
	}

	var request struct {
		Country string `json:"country"`
		City    string `json:"city"`
		Type    string `json:"type"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	placeDetail, err := services.GetPlaceDetail(placeName, request.Type, request.Country, request.City)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch place details"})
		return
	}

	c.JSON(http.StatusOK, placeDetail)
}

type Place struct {
	Image    string `json:"image"`
	Landmark string `json:"landmark"`
	RoadName string `json:"roadName"`
	Time     string `json:"time"`
	Town     string `json:"town"`
	Type     string `json:"type"`
}

type TimelineResponse struct {
	Budget   string             `json:"budget"`
	Country  string             `json:"country"`
	Town     string             `json:"town"`
	Title    string             `json:"title"`
	Timeline map[string][]Place `json:"timeline"`
}

func convertPlaces(places []struct {
	Name      string `json:"name"`
	TimeOfDay string `json:"timeOfDay"`
}) []services.PlaceTimeline {
	var converted []services.PlaceTimeline
	for _, place := range places {
		converted = append(converted, services.PlaceTimeline{
			Name:      place.Name,
			TimeOfDay: place.TimeOfDay,
		})
	}
	return converted
}

func Timeline(c *gin.Context) {
	var request struct {
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

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	response, err := services.GetTimelineFromGemini(
		request.Accomodation,
		request.Town,
		request.Country,
		request.StartDate,
		request.EndDate,
		convertPlaces(request.Places),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get response from Gemini"})
		return
	}

	timeline, ok := response["timeline"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response format from Gemini API"})
		return
	}

	formattedResponse := gin.H{
		"budget":   response["budget"],
		"country":  response["country"],
		"town":     response["town"],
		"title":    response["title"],
		"timeline": timeline,
	}

	c.JSON(http.StatusOK, formattedResponse)
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
		return time.Now() // Jika gagal parsing, gunakan waktu saat ini
	}
	return parsedTime
}

func SelectPlace(c *gin.Context) {
	var request SelectPlaceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	history := models.History{
		Country:     request.Country,
		Town:        request.Town,
		StartDate:   parseDate(request.StartDate),
		EndDate:     parseDate(request.EndDate),
		Description: request.Title,
		Image:       models.StringArray{},
	}

	history.SelectedAccomodation = []models.SelectedAccomodation{
		{
			Name:  request.Accomodation,
			Image: models.StringArray{},
		},
	}
	for _, timelineDetails := range request.Timelines {
		for _, detail := range timelineDetails {
			history.SelectedPlaces = append(history.SelectedPlaces, models.SelectedPlace{
				PlaceToVisit: detail.Landmark,
				Town:         detail.Town,
				Image:        models.StringArray{detail.Image},
			})
		}
	}

	if err := database.DB.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save history"})
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

	c.JSON(http.StatusOK, response)
}
