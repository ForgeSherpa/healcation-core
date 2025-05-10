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

type AIService interface {
	Search(query string) ([]PlaceSearch, error)
	GetPlaces(preferences []string, country, town string) (map[string]interface{}, error)
	GetPlaceDetail(name, placeType, country, city string) (PlaceDetail, error)
	GetTimeline(accommodation, town, country, startDate, endDate string, places []struct {
		Name      string `json:"name"`
		TimeOfDay string `json:"timeOfDay"`
	}) (map[string]interface{}, error)
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
