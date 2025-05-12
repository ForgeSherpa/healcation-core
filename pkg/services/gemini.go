package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"healcationBackend/pkg/config"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type GeminiService struct {
	apiKey string
}

func NewGeminiService() AIService {
	return GeminiService{
		apiKey: config.GeminiAPIKey,
	}
}

func removeMarkdownCodeBlock(input string) string {
	re := regexp.MustCompile("(?s)```json\\n(.*?)\\n```")
	matches := re.FindStringSubmatch(input)

	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return strings.TrimSpace(input)
}
func cleanJSONResponse(response string) string {
	re := regexp.MustCompile("(?s)```json(.*?)```")
	cleaned := re.ReplaceAllString(response, "$1")

	cleaned = strings.TrimSpace(cleaned)
	return cleaned
}

// Fitur Search
type GeminiResponseSearch struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (s GeminiService) Search(query string) ([]PlaceSearch, error) {
	apiKey := config.GeminiAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("API Key tidak ditemukan")
	}

	apiURL := config.GeminiAPIKey

	prompt := fmt.Sprintf(`Cari informasi tentang destinasi "%s" dan berikan respons dalam format JSON seperti berikut:
	{
	  "results": [
	    {
	      "country": "Nama negara",
	      "town": "Nama kota atau daerah"
	    }
	  ]
	}
	Hanya kembalikan JSON di atas tanpa teks tambahan.`, query)

	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role":  "user",
				"parts": []map[string]string{{"text": prompt}},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geminiResp GeminiResponseSearch
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("respon dari Gemini kosong atau tidak sesuai format")
	}

	rawJSON := geminiResp.Candidates[0].Content.Parts[0].Text
	cleanedJSON := cleanJSONResponse(rawJSON)

	var result struct {
		Results []PlaceSearch `json:"results"`
	}
	if err := json.Unmarshal([]byte(cleanedJSON), &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

// Fitur GetPlaces
type PlaceGetPlaces struct {
	Description string   `json:"description"`
	Image       []string `json:"image"` // Mengubah image menjadi array
	Name        string   `json:"name"`
	Town        string   `json:"town"`
	Type        string   `json:"type"`
}

type AccommodationGetPlaces struct {
	Image []string `json:"image"` // Mengubah image menjadi array
	Name  string   `json:"name"`
}

type GeminiResponseGetPlaces struct {
	Accomodations []AccommodationGetPlaces `json:"accomodations"`
	Places        []PlaceGetPlaces         `json:"places"`
}

type GeminiResponseSearchGetPlaces struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (s GeminiService) GetPlaces(preferences []string, country, town string) (map[string]interface{}, error) {
	apiKey := config.GeminiAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("API Key tidak ditemukan")
	}

	apiURL := config.GeminiAPIKey

	prompt := fmt.Sprintf(`Berikan daftar tempat wisata dan akomodasi di %s, %s berdasarkan preferensi berikut: %v.
Harap berikan respons dalam format JSON dengan struktur berikut:

{
  "preferences": %v,
  "country": "%s",
  "town": "%s",
  "accomodations": [
    {
      "name": "Nama akomodasi"
    }
  ],
  "places": [
    {
      "description": "Deskripsi singkat tentang tempat wisata",
      "name": "Nama tempat wisata",
      "town": "%s",
      "type": "Jenis tempat wisata (contoh: Museum, Landmark, District, dll)"
    }
  ]
}
Hanya kembalikan JSON di atas tanpa teks tambahan.`, town, country, preferences, preferences, country, town, town)

	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role":  "user",
				"parts": []map[string]string{{"text": prompt}},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var geminiResponse struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResponse); err != nil {
		return nil, fmt.Errorf("gagal parsing JSON response: %v", err)
	}

	if len(geminiResponse.Candidates) == 0 || len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("respon dari Gemini kosong atau tidak sesuai format")
	}

	rawJSON := removeMarkdownCodeBlock(geminiResponse.Candidates[0].Content.Parts[0].Text)

	var result struct {
		Places        []PlaceGetPlaces         `json:"places"`
		Accomodations []AccommodationGetPlaces `json:"accomodations"`
	}

	if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
		return nil, fmt.Errorf("gagal parsing JSON teks: %v", err)
	}

	for i := range result.Places {
		name := result.Places[i].Name
		imageURLs, err := GetGoogleImagesPlaces(name)
		if err == nil {
			result.Places[i].Image = imageURLs
		}
	}

	for i := range result.Accomodations {
		name := result.Accomodations[i].Name
		imageURLs, err := GetGoogleImages(name)
		if err == nil {
			result.Accomodations[i].Image = imageURLs
		}
	}

	response := map[string]interface{}{
		"accomodations": result.Accomodations,
		"places":        result.Places,
	}

	return response, nil
}

func (s GeminiService) GetPlaceDetail(name, placeType, country, city string) (PlaceDetail, error) {
	apiKey := config.GeminiAPIKey
	if apiKey == "" {
		return PlaceDetail{}, fmt.Errorf("API Key tidak ditemukan")
	}

	apiURL := config.GeminiAPIKey // TODO: kayaknya ini salah, harusnya ke endpoint Gemini

	typeInfo := ""
	if placeType == "accommodation" {
		typeInfo = "Tempat ini adalah akomodasi (hotel atau penginapan)."
	}

	prompt := fmt.Sprintf(`Berikan informasi tentang tempat bernama "%s" di kota %s, %s. %s
Harap kembalikan data dalam format JSON sebagai berikut:

{
  "name": "%s",
  "description": "Deskripsi singkat"
}

Hanya kembalikan JSON di atas tanpa teks tambahan.`, name, city, country, typeInfo, name)

	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role":  "user",
				"parts": []map[string]string{{"text": prompt}},
			},
		},
	})
	if err != nil {
		return PlaceDetail{}, err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return PlaceDetail{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PlaceDetail{}, err
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return PlaceDetail{}, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return PlaceDetail{}, fmt.Errorf("respon dari Gemini kosong atau tidak sesuai format")
	}

	rawJSON := removeMarkdownCodeBlock(geminiResp.Candidates[0].Content.Parts[0].Text)

	var placeDetail PlaceDetail
	if err := json.Unmarshal([]byte(rawJSON), &placeDetail); err != nil {
		return PlaceDetail{}, err
	}

	imageURLs, err := GetGoogleImages(name)
	if err == nil {
		placeDetail.Image = imageURLs[0]
	}

	return placeDetail, nil
}

// fitur timeline
type PlaceTimeline struct {
	Name      string `json:"name"`
	TimeOfDay string `json:"timeOfDay"`
}

type Place struct {
	Image       string `json:"image"`
	Landmark    string `json:"landmark"`
	RoadName    string `json:"roadName"`
	Time        string `json:"time"`
	Town        string `json:"town"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (s GeminiService) GetTimeline(accommodation, town, country, startDate, endDate string, places []struct {
	Name      string `json:"name"`
	TimeOfDay string `json:"timeOfDay"`
}) (map[string]interface{}, error) {
	apiKey := config.GeminiAPIKey
	if apiKey == "" {
		return nil, errors.New("API Key tidak ditemukan")
	}

	apiURL := config.GeminiAPIKey // TODO: kayaknya ini salah, harusnya ke endpoint Gemini

	prompt := fmt.Sprintf(`Buatkan rencana perjalanan dari %s, %s pada tanggal %s hingga %s berdasarkan tempat berikut: %v.
Harap berikan respons dalam format JSON dengan struktur berikut:

{
  "budget": "estimasi dalam IDR ; contoh 1.000.000",
  "country": "%s",
  "town": "%s",
  "title": "Gemini Generated",
  "timeline": {
    "YYYY-MM-DD": [
      {
        "landmark": "Nama tempat",
        "roadName": "Nama jalan (jangan kosong atau N/A, selalu isi dengan jalan yang relevan)",
        "time": "Waktu kunjungan",
        "town": "%s",
        "type": "Hotel ; Staycation ; Food ; Event ; Local Event ; Tourist Attraction ; Museum ; Historical Site ; History Site ; Cultural Site",
		"description": "Deskripsi lengkap mengenai tempat wisata, dalam 3 kalimat, isinya semua berupa pesona tempat wisatra saja ; contoh Nikmati pesiar sungai Seine di malam hari. Saksikan landmark-landmark kota diterangi dengan indah. Ini adalah cara romantis dan santai untuk melihat Paris dari perspektif yang berbeda.",
      }
    ]
  }
}
Hanya kembalikan JSON di atas tanpa teks tambahan.`, town, country, startDate, endDate, places, country, town, town)

	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role":  "user",
				"parts": []map[string]string{{"text": prompt}},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geminiResponse struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResponse); err != nil {
		return nil, fmt.Errorf("gagal parsing JSON response: %v", err)
	}

	if len(geminiResponse.Candidates) == 0 || len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("respon dari Gemini kosong atau tidak sesuai format")
	}

	rawJSON := removeMarkdownCodeBlock(geminiResponse.Candidates[0].Content.Parts[0].Text)

	var result struct {
		Budget   string             `json:"budget"`
		Country  string             `json:"country"`
		Town     string             `json:"town"`
		Title    string             `json:"title"`
		Timeline map[string][]Place `json:"timeline"`
	}

	if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
		return nil, fmt.Errorf("gagal parsing JSON teks: %v", err)
	}

	for date, places := range result.Timeline {
		for i := range places {
			name := places[i].Landmark
			imageURLs, err := GetGoogleImages(name)
			if err == nil && len(imageURLs) > 0 {
				places[i].Image = imageURLs[0]
			}
		}
		result.Timeline[date] = places
	}

	response := map[string]interface{}{
		"budget":   result.Budget,
		"country":  result.Country,
		"town":     result.Town,
		"title":    result.Title,
		"timeline": result.Timeline,
	}

	return response, nil
}
