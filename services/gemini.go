package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

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

	body, err := ioutil.ReadAll(resp.Body)
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
		fmt.Println("❌ API Key tidak ditemukan di environment!")
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
      "image": "URL gambar akomodasi",
      "name": "Nama akomodasi"
    }
  ],
  "places": [
    {
      "description": "Deskripsi singkat tentang tempat wisata",
      "image": "URL gambar tempat wisata",
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
		fmt.Println("❌ Gagal membuat JSON request:", err)
		return nil, err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("❌ Gagal menghubungi Gemini API:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("❌ Gagal membaca response dari Gemini:", err)
		return nil, err
	}

	fmt.Println("📥 Respon dari Gemini API:", string(body))

	var geminiResp GeminiResponseSearchGetPlaces
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		fmt.Println("❌ Gagal parsing JSON response:", err)
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		fmt.Println("❌ Respon dari Gemini kosong atau tidak sesuai format")
		return nil, fmt.Errorf("respon dari Gemini kosong atau tidak sesuai format")
	}

	rawJSON := geminiResp.Candidates[0].Content.Parts[0].Text
	cleanedJSON := cleanJSONResponse(rawJSON)

	fmt.Println("✅ JSON setelah dibersihkan:", cleanedJSON)

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleanedJSON), &result); err != nil {
		fmt.Println("❌ Gagal parsing JSON setelah pembersihan:", err)
		return nil, err
	}

	delete(result, "preferences")
	delete(result, "town")
	delete(result, "country")

	return result, nil
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
  "image": "URL gambar tempat",
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

	body, err := ioutil.ReadAll(resp.Body)
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

	rawJSON := geminiResp.Candidates[0].Content.Parts[0].Text
	cleanedJSON := cleanJSONResponse(rawJSON)

	var placeDetail PlaceDetail
	if err := json.Unmarshal([]byte(cleanedJSON), &placeDetail); err != nil {
		return PlaceDetail{}, err
	}

	return placeDetail, nil
}

// Fitur Timeline
type PlaceTimeline struct {
	Name      string `json:"name"`
	TimeOfDay string `json:"timeOfDay"`
}

func GetTimelineFromGemini(accomodation, town, country, startDate, endDate string, places []PlaceTimeline) (map[string]interface{}, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ API Key tidak ditemukan di environment!")
		return nil, fmt.Errorf("API Key tidak ditemukan")
	}

	apiURL := os.Getenv("GEMINI_API_KEY")

	placesList := ""
	for _, place := range places {
		placesList += fmt.Sprintf("{'name': '%s', 'timeOfDay': '%s'}, ", place.Name, place.TimeOfDay)
	}

	prompt := fmt.Sprintf(`Buatkan itinerary perjalanan untuk akomodasi %s di kota %s, negara %s, mulai dari tanggal %s sampai %s.
Tempat yang akan dikunjungi adalah sebagai berikut: %s

Format respons yang diharapkan:
{
    "budget": "min - max",
    "country": "%s",
    "town": "%s",
    "title": "Gemini Generated",
    "timeline": {
        "YYYY-MM-DD": [
            {
                "image": "URL gambar tempat",
                "landmark": "Nama tempat",
                "roadName": "Nama jalan",
                "time": "HH:MM",
                "town": "Nama kota",
                "type": "Jenis tempat (contoh: Tourist Attraction, Park, Museum)"
            }
        ]
    }
}
Hanya kembalikan JSON di atas tanpa teks tambahan.`, accomodation, town, country, startDate, endDate, placesList, country, town)

	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role":  "user",
				"parts": []map[string]string{{"text": prompt}},
			},
		},
	})
	if err != nil {
		fmt.Println("❌ Gagal membuat JSON request:", err)
		return nil, err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("❌ Gagal menghubungi Gemini API:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("❌ Gagal membaca response dari Gemini:", err)
		return nil, err
	}
	fmt.Println("🔍 Response dari Gemini API:", string(body)) // Debugging output

	var geminiResponse map[string]interface{}
	if err := json.Unmarshal(body, &geminiResponse); err != nil {
		fmt.Println("❌ Gagal parsing JSON dari Gemini:", err)
		return nil, err
	}

	candidates, ok := geminiResponse["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		fmt.Println("❌ Tidak ada kandidat dalam response Gemini API")
		return nil, fmt.Errorf("response tidak memiliki kandidat yang valid")
	}

	content, ok := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Response tidak memiliki format konten yang diharapkan")
		return nil, fmt.Errorf("format konten tidak valid")
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		fmt.Println("❌ Bagian teks tidak ditemukan dalam respons")
		return nil, fmt.Errorf("bagian teks kosong")
	}

	textResponse, ok := parts[0].(map[string]interface{})["text"].(string)
	if !ok {
		fmt.Println("❌ Respons tidak memiliki teks yang valid")
		return nil, fmt.Errorf("format teks tidak valid")
	}

	cleanedText := strings.TrimPrefix(textResponse, "```json\n")
	cleanedText = strings.TrimSuffix(cleanedText, "\n```")

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleanedText), &result); err != nil {
		fmt.Println("❌ Gagal parsing JSON hasil:", err)
		return nil, err
	}

	return result, nil
}
