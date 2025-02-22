package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"healcationBackend/models"
)

// Struktur untuk menerima respons dari Gemini API
type GeminiResponse struct {
	PlaceRecommendations        []models.PlaceRecommendation        `json:"placeRecommendation"`
	AccomodationRecommendations []models.AccomodationRecommendation `json:"accomodationRecommendation"`
}

func FetchFromGeminiAPI(city string) (*GeminiResponse, error) {
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"

	// Ambil OAuth 2.0 token dari gcloud
	tokenCmd := exec.Command("gcloud", "auth", "application-default", "print-access-token")
	tokenBytes, err := tokenCmd.Output()
	if err != nil {
		log.Println("Error fetching OAuth token:", err)
		return nil, err
	}
	oauthToken := strings.TrimSpace(string(tokenBytes))

	// Data yang dikirim ke Gemini API (sesuai format)
	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]string{
			{"parts": city},
		},
	})
	if err != nil {
		return nil, err
	}

	// Buat request ke API Gemini
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+oauthToken) // Pakai OAuth Token

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error client.Do:", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Println("Status Code:", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return nil, err
	}
	log.Println("Response Body:", string(body))

	// Jika status bukan 200, return error
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch data from Gemini API: " + string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		log.Println("Error unmarshal:", err)
		return nil, err
	}

	return &geminiResp, nil
}
