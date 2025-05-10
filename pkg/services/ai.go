package services

import "healcationBackend/pkg/config"

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

func NewAIService() AIService {
	if config.IsGeminiEnabled {
		return NewGeminiService()
	}

	return NewGeminiMockService()
}
