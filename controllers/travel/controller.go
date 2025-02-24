package travel

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlaces(c *gin.Context) {
	response := gin.H{
		"accomodations": []gin.H{
			{
				"name":  "Luxury Hotel Paris",
				"image": "https://example.com/luxury-hotel.jpg",
			},
		},
		"places": []gin.H{
			{
				"name":        "Eiffel Tower",
				"image":       "https://example.com/eiffel-tower.jpg",
				"description": "A famous landmark in Paris, known for its stunning views.",
				"town":        "Paris",
				"type":        "Landmark",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

func GetPlaceDetail(c *gin.Context) {
	placeName := c.Param("name")

	response := gin.H{
		"name":        placeName,
		"image":       "https://example.com/img.jpg",
		"description": "Contoh deskripsi tempat wisata atau akomodasi.",
	}

	c.JSON(http.StatusOK, response)
}

func Timeline(c *gin.Context) {
	response := gin.H{
		"budget":  "100 - 500 USD",
		"country": "Indonesia",
		"town":    "Jakarta",
		"title":   "Gemini Generated",
		"timeline": map[string][]map[string]string{
			"2024-04-01": {
				{
					"image":    "barelang1.jpg",
					"landmark": "Barelang Bridge",
					"roadName": "Jl. Trans Barelang",
					"time":     "10:00",
					"town":     "Batam",
					"type":     "Tourist Attraction",
				},
				{
					"image":    "barelang1.jpg",
					"landmark": "Batam Botanical Garden",
					"roadName": "Jl. Engku Putri",
					"time":     "14:00",
					"town":     "Batam",
					"type":     "Park",
				},
			},
		},
	}

	c.JSON(http.StatusOK, response)
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

func SelectPlace(c *gin.Context) {
	var request SelectPlaceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
