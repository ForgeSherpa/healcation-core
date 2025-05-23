package services

import (
	"encoding/json"
	"errors"
	"healcationBackend/pkg/config"
	"net/http"
)

type PlaceSearch struct {
	Country string `json:"country"`
	Town    string `json:"town"`
}

type PlaceDetail struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type SelectedPlace struct {
	TimeOfDay string   `json:"timeOfDay"`
	Places    []string `json:"places"`
}

type PlaceGetPlaces struct {
	Description string   `json:"description"`
	Image       []string `json:"image"`
	Name        string   `json:"name"`
	Town        string   `json:"town"`
	Type        string   `json:"type"`
}

type AccommodationGetPlaces struct {
	Image []string `json:"image"`
	Name  string   `json:"name"`
	Town  string   `json:"town"`
}

type GeminiResponseGetPlaces struct {
	Accomodations []AccommodationGetPlaces `json:"accomodations"`
	Places        []PlaceGetPlaces         `json:"places"`
}

type PlaceVisitedDetail struct {
	Type     string   `json:"type"`
	Landmark string   `json:"landmark"`
	RoadName string   `json:"roadName"`
	Town     string   `json:"town"`
	Time     string   `json:"time"`
	Image    []string `json:"image"`
}

type DailyVisit struct {
	Date string               `json:"date"`
	Data []PlaceVisitedDetail `json:"data"`
}

type TimelineResponse struct {
	Budget       string       `json:"budget"`
	Town         string       `json:"town"`
	Country      string       `json:"country"`
	StartDate    string       `json:"startDate"`
	EndDate      string       `json:"endDate"`
	PlaceVisited []DailyVisit `json:"placeVisited"`
}

type LandmarkDetail struct {
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

type AIService interface {
	Search(query string) ([]PlaceSearch, error)
	GetPlaces(preferences []string, country, town string) (map[string]interface{}, error)
	GetTimeline(
		accommodation, town, country, startDate, endDate string,
		places []SelectedPlace,
	) (*TimelineResponse, error)

	GetPlaceDetail(placeType, landmark, town string) (map[string]interface{}, error)
}

var ErrGeminiUnavailable = errors.New("gemini service is unavailable")

func HandleGeminiUnavailable(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrGeminiUnavailable) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": ErrGeminiUnavailable.Error(),
		})
	}
}

func NewAIService() (AIService, error) {
	if config.IsStaging {
		return NewGeminiMockService(), nil
	}

	if config.IsGeminiEnabled {
		return NewGeminiService(), nil
	}

	return nil, ErrGeminiUnavailable
}
