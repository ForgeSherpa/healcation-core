package database

import (
	"healcationBackend/models"
	"time"
)

func SeedPopular() {
	populars := []models.Popular{
		{ID: "1", Town: "Bali", Country: "Indonesia", Image: "https://example.com/bali.jpg"},
		{ID: "2", Town: "Yogyakarta", Country: "Indonesia", Image: "https://example.com/jogja.jpg"},
		{ID: "3", Town: "Lombok", Country: "Indonesia", Image: "https://example.com/lombok.jpg"},
	}

	DB.Create(&populars)
}

func SeedHistory() {
	histories := []models.History{
		{
			ID:        "1",
			Town:      "Bali",
			Country:   "Indonesia",
			StartDate: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2023, time.January, 7, 0, 0, 0, 0, time.UTC),
			Image:     "https://example.com/bali.jpg",
			UserID:    "1",
		},
		//{
		//	ID:        "2",
		//	Town:      "Yogyakarta",
		//	Country:   "Indonesia",
		//	StartDate: time.Date(2023, time.February, 10, 0, 0, 0, 0, time.UTC),
		//	EndDate:   time.Date(2023, time.February, 15, 0, 0, 0, 0, time.UTC),
		//	Image:     "https://example.com/jogja.jpg",
		//	UserID:    "1",
		//},
		//{
		//	ID:        "3",
		//	Town:      "Lombok",
		//	Country:   "Indonesia",
		//	StartDate: time.Date(2023, time.March, 5, 0, 0, 0, 0, time.UTC),
		//	EndDate:   time.Date(2023, time.March, 12, 0, 0, 0, 0, time.UTC),
		//	Image:     "https://example.com/lombok.jpg",
		//	UserID:    "1",
		//},
	}
	DB.Create(&histories)
}

func SeedPreference() {
	preferences := []models.Preferences{
		{ID: "1", StartDate: "2025-01-01", EndDate: "2025-01-10"},
		{ID: "2", StartDate: "2025-02-01", EndDate: "2025-02-10"},
	}
	DB.Create(&preferences)

	preferenceTypes := []models.Preference{
		{ID: "1", Type: "attraction", Image: "https://example.com/attraction.jpg"},
		{ID: "2", Type: "foodies", Image: "https://example.com/foodies.jpg"},
		{ID: "3", Type: "staycation", Image: "https://example.com/staycation.jpg"},
		{ID: "4", Type: "shopping", Image: "https://example.com/shopping.jpg"},
		{ID: "5", Type: "historical site", Image: "https://example.com/historical.jpg"},
		{ID: "6", Type: "outdoor activity", Image: "https://example.com/outdoor.jpg"},
	}
	DB.Create(&preferenceTypes)

	// Link Preferences and Preference Types (join table with additional fields)
	joinData := []models.PreferenceLink{
		{PreferencesID: "1", PreferenceID: "1", Type: "attraction", Image: "https://example.com/attraction.jpg"},
		{PreferencesID: "1", PreferenceID: "2", Type: "foodies", Image: "https://example.com/foodies.jpg"},
		{PreferencesID: "2", PreferenceID: "3", Type: "staycation", Image: "https://example.com/staycation.jpg"},
		{PreferencesID: "2", PreferenceID: "4", Type: "shopping", Image: "https://example.com/shopping.jpg"},
	}
	DB.Create(&joinData)
}

func SeedSelect_Place() {
	times := []models.Time{
		{ID: "1", TimeOfDay: "morning"},
		{ID: "2", TimeOfDay: "afternoon"},
		{ID: "3", TimeOfDay: "night"},
	}
	for _, timez := range times {
		var existing models.Time
		DB.Where("TimeOfDay = ?", timez.TimeOfDay).First(&existing)
		if existing.ID == "" {
			DB.Create(&timez)
		}
	}

	selectPlaces := []models.SelectPlace{
		{
			ID:          "1",
			City:        "Paris",
			Country:     "France",
			Description: "Eiffel Tower in the morning",
			PlaceRecommendation: []models.PlaceRecommendation{
				{ID: "1", PlaceToVisit: "Eiffel Tower"},
			},
			SelectedPlace: []models.Time{
				{ID: "1", TimeOfDay: "morning"},
			},
		},
		{
			ID:          "2",
			City:        "Paris",
			Country:     "France",
			Description: "Louvre Museum in the morning",
			PlaceRecommendation: []models.PlaceRecommendation{
				{ID: "2", PlaceToVisit: "Louvre Museum"},
			},
			SelectedPlace: []models.Time{
				{ID: "1", TimeOfDay: "morning"},
			},
		},
		{
			ID:          "3",
			City:        "Paris",
			Country:     "France",
			Description: "Eiffel Tower in the afternoon",
			PlaceRecommendation: []models.PlaceRecommendation{
				{ID: "3", PlaceToVisit: "Eiffel Tower"},
			},
			SelectedPlace: []models.Time{
				{ID: "2", TimeOfDay: "afternoon"},
			},
		},
		{
			ID:          "4",
			City:        "Paris",
			Country:     "France",
			Description: "Louvre Museum in the afternoon",
			PlaceRecommendation: []models.PlaceRecommendation{
				{ID: "4", PlaceToVisit: "Louvre Museum"},
			},
			SelectedPlace: []models.Time{
				{ID: "2", TimeOfDay: "afternoon"},
			},
		},
		{
			ID:          "5",
			City:        "Paris",
			Country:     "France",
			Description: "Eiffel Tower at night",
			PlaceRecommendation: []models.PlaceRecommendation{
				{ID: "5", PlaceToVisit: "Eiffel Tower"},
			},
			SelectedPlace: []models.Time{
				{ID: "3", TimeOfDay: "night"},
			},
		},
		{
			ID:          "6",
			City:        "Paris",
			Country:     "France",
			Description: "Louvre Museum at night",
			PlaceRecommendation: []models.PlaceRecommendation{
				{ID: "6", PlaceToVisit: "Louvre Museum"},
			},
			SelectedPlace: []models.Time{
				{ID: "3", TimeOfDay: "night"},
			},
		},
	}
	DB.Create(&selectPlaces)
	//for _, selectPlace := range selectPlaces {
	//	var existing models.SelectPlace
	//	DB.Where("ID = ?", selectPlace.ID).First(&existing)
	//	if existing.ID == "" {
	//		DB.Create(&selectPlace)
	//	}
	//}
}

func SeedTimelines() {
	// Data timelines
	timelines := []models.Timeline{
		{
			ID:        "1",
			Town:      "Batam",
			Country:   "Indonesia",
			Budget:    "5000",
			StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        "2",
			Town:      "Jakarta",
			Country:   "Indonesia",
			Budget:    "8000",
			StartDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC),
		},
	}

	placeVisited := []models.PlaceVisited{
		{
			TimelineID: "1",
			Type:       "Tourist Attraction",
			Landmark:   "Barelang Bridge",
			RoadName:   "Jl. Trans Barelang",
			Town:       "Batam",
			Time:       "10:00",
			Images:     `["barelang1.jpg", "barelang2.jpg"]`,
		},
		{
			TimelineID: "1",
			Type:       "Park",
			Landmark:   "Batam Botanical Garden",
			RoadName:   "Jl. Engku Putri",
			Town:       "Batam",
			Time:       "14:00",
			Images:     `["garden1.jpg", "garden2.jpg"]`,
		},
		{
			TimelineID: "2",
			Type:       "Museum",
			Landmark:   "National Museum",
			RoadName:   "Jl. Medan Merdeka Barat",
			Town:       "Jakarta",
			Time:       "11:00",
			Images:     `["museum1.jpg", "museum2.jpg"]`,
		},
		{
			TimelineID: "2",
			Type:       "Shopping Mall",
			Landmark:   "Grand Indonesia",
			RoadName:   "Jl. M.H. Thamrin",
			Town:       "Jakarta",
			Time:       "15:00",
			Images:     `["mall1.jpg", "mall2.jpg"]`,
		},
	}

	DB.Create(&timelines)
	DB.Create(&placeVisited)
}
