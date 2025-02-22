package travel

import (
	"healcationBackend/database"
	"healcationBackend/models"
	"net/http"

	// "healcationBackend/services"
	"github.com/gin-gonic/gin"
)

func GetPreferences(c *gin.Context) {
	var preferenceLinks []models.PreferenceLink
	if err := database.DB.Find(&preferenceLinks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve preferences"})
		return
	}

	response := make([]map[string]interface{}, len(preferenceLinks))
	for i, link := range preferenceLinks {
		response[i] = map[string]interface{}{
			"preferencesId": link.PreferencesID,
			"preferenceId":  link.PreferenceID,
			"type":          link.Type,
			"image":         link.Image,
		}
	}

	c.JSON(http.StatusOK, gin.H{"preferences": response})
}

// func GetSelectPlaces(c *gin.Context) {
// 	var selectPlaces []models.SelectPlace
// 	if err := database.DB.Preload("PlaceRecommendation").Preload("AccomodationRecommendation").Find(&selectPlaces).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve select places"})
// 		return
// 	}

// 	// Looping untuk mengambil data dari Gemini API berdasarkan kota
// 	for i, place := range selectPlaces {
// 		geminiData, err := services.FetchFromGeminiAPI(place.City)
// 		if err == nil && geminiData != nil {
// 			// Update data di database
// 			database.DB.Model(&place).Association("PlaceRecommendation").Replace(geminiData.PlaceRecommendations)
// 			database.DB.Model(&place).Association("AccomodationRecommendation").Replace(geminiData.AccomodationRecommendations)

// 			// Update response dengan data terbaru
// 			selectPlaces[i].PlaceRecommendation = geminiData.PlaceRecommendations
// 			selectPlaces[i].AccomodationRecommendation = geminiData.AccomodationRecommendations
// 		}
// 	}

// 	// Buat respons JSON
// 	response := make([]map[string]interface{}, len(selectPlaces))
// 	for i, place := range selectPlaces {
// 		response[i] = map[string]interface{}{
// 			"id":                         place.ID,
// 			"city":                       place.City,
// 			"country":                    place.Country,
// 			"description":                place.Description,
// 			"placeRecommendation":        place.PlaceRecommendation,
// 			"accomodationRecommendation": place.AccomodationRecommendation,
// 			"selectedPlace":              place.SelectedPlace,
// 			"selectedAccomodation":       place.SelectedAccomodation,
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{"selectPlaces": response})
// }

// GetSelectPlaces mengambil daftar SelectPlace beserta relasinya
func GetSelectPlaces(c *gin.Context) {
	var selectPlaces []models.SelectPlace

	// Menggunakan Preload untuk mengambil relasi
	if err := database.DB.
		Preload("PlaceRecommendation").
		Preload("AccomodationRecommendation").
		Preload("SelectedPlace.Places").
		Find(&selectPlaces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve select places"})
		return
	}

	// Membentuk respons JSON
	response := make([]map[string]interface{}, len(selectPlaces))
	for i, place := range selectPlaces {
		// Format data PlaceRecommendation
		placeRecommendations := make([]map[string]interface{}, len(place.PlaceRecommendation))
		for j, rec := range place.PlaceRecommendation {
			placeRecommendations[j] = map[string]interface{}{
				"id":           rec.ID,
				"placeToVisit": rec.PlaceToVisit,
				"town":         rec.Town,
				"image":        rec.Image,
			}
		}

		// Format data AccomodationRecommendation
		accommodationRecommendations := make([]map[string]interface{}, len(place.AccomodationRecommendation))
		for j, rec := range place.AccomodationRecommendation {
			accommodationRecommendations[j] = map[string]interface{}{
				"id":    rec.ID,
				"name":  rec.Name,
				"town":  rec.Town,
				"image": rec.Image,
			}
		}

		// Format data SelectedPlace
		selectedPlaces := make([]map[string]interface{}, len(place.SelectedPlace))
		for j, timeSlot := range place.SelectedPlace {
			places := make([]map[string]interface{}, len(timeSlot.Places))
			for k, selectedPlace := range timeSlot.Places {
				places[k] = map[string]interface{}{
					"id":           selectedPlace.ID,
					"placeToVisit": selectedPlace.PlaceToVisit,
					"town":         selectedPlace.Town,
				}
			}

			selectedPlaces[j] = map[string]interface{}{
				"id":        timeSlot.ID,
				"timeOfDay": timeSlot.TimeOfDay,
				"places":    places,
			}
		}

		response[i] = map[string]interface{}{
			"id":                         place.ID,
			"city":                       place.City,
			"country":                    place.Country,
			"description":                place.Description,
			"placeRecommendation":        placeRecommendations,
			"accomodationRecommendation": accommodationRecommendations,
			"selectedPlace":              selectedPlaces,
			"selectedAccomodation":       place.SelectedAccomodation,
		}
	}

	c.JSON(http.StatusOK, gin.H{"selectPlaces": response})
}
