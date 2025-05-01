package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func GetGoogleImages(query string) ([]string, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	cx := os.Getenv("GOOGLE_API_CX")

	if apiKey == "" || cx == "" {
		return nil, fmt.Errorf("google API Key atau CX tidak ditemukan")
	}

	apiURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&key=%s&cx=%s&searchType=image&num=1",
		url.QueryEscape(query), apiKey, cx)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gagal mendapatkan gambar dari Google: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleResp struct {
		Items []struct {
			Link string `json:"link"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &googleResp); err != nil {
		return nil, err
	}

	if len(googleResp.Items) == 0 {
		return nil, fmt.Errorf("tidak ada gambar ditemukan untuk query: %s", query)
	}

	return []string{googleResp.Items[0].Link}, nil
}

func GetGoogleImagesPlaces(query string) ([]string, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	cx := os.Getenv("GOOGLE_API_CX")

	if apiKey == "" || cx == "" {
		return nil, fmt.Errorf("google API Key atau CX tidak ditemukan")
	}

	apiURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&key=%s&cx=%s&searchType=image&num=2",
		url.QueryEscape(query), apiKey, cx)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gagal mendapatkan gambar dari Google: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleResp struct {
		Items []struct {
			Link string `json:"link"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &googleResp); err != nil {
		return nil, err
	}

	var imageURLs []string
	for i := 0; i < len(googleResp.Items) && i < 2; i++ {
		imageURLs = append(imageURLs, googleResp.Items[i].Link)
	}

	if len(imageURLs) == 0 {
		return nil, fmt.Errorf("tidak ada gambar ditemukan untuk query: %s", query)
	}

	return imageURLs, nil
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
type PlaceSearch struct {
	Country string `json:"country"`
	Town    string `json:"town"`
}

type GeminiResponseSearch struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func SearchGemini(query string) ([]PlaceSearch, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API Key tidak ditemukan")
	}

	apiURL := os.Getenv("GEMINI_API_KEY")

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

func GetPlacesFromGemini(preferences []string, country, town string) (map[string]interface{}, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API Key tidak ditemukan")
	}

	apiURL := os.Getenv("GEMINI_API_KEY")

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

// fitur GetPlaceDetail
type PlaceDetail struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

func GetPlaceDetail(name, placeType, country, city string) (PlaceDetail, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return PlaceDetail{}, fmt.Errorf("API Key tidak ditemukan")
	}

	apiURL := os.Getenv("GEMINI_API_KEY")

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
	Image    string `json:"image"`
	Landmark string `json:"landmark"`
	RoadName string `json:"roadName"`
	Time     string `json:"time"`
	Town     string `json:"town"`
	Type     string `json:"type"`
}

func GetTimelineFromGemini(accommodation, town, country, startDate, endDate string, places []struct {
	Name      string `json:"name"`
	TimeOfDay string `json:"timeOfDay"`
}) (map[string]interface{}, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("API Key tidak ditemukan")
	}

	apiURL := os.Getenv("GEMINI_API_KEY")

	prompt := fmt.Sprintf(`Buatkan rencana perjalanan dari %s, %s pada tanggal %s hingga %s berdasarkan tempat berikut: %v.
Harap berikan respons dalam format JSON dengan struktur berikut:

{
  "budget": "min - max dalam IDR",
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
        "type": "Jenis tempat (Museum, Landmark, dll)"
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
