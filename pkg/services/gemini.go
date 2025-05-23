package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"healcationBackend/pkg/config"
	"regexp"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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
func (s GeminiService) Search(query string) ([]PlaceSearch, error) {
	apiKey := config.GeminiAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("API Key tidak ditemukan")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")

	responseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"results": {
				Type:        genai.TypeArray,
				Description: "Daftar hasil pencarian destinasi.",
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"country": {Type: genai.TypeString, Description: "Nama negara dari destinasi."},
						"town":    {Type: genai.TypeString, Description: "Nama kota atau daerah dari destinasi."},
					},
					Required: []string{"country", "town"},
				},
			},
		},
		Required: []string{"results"},
	}

	model.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	}

	instructionLanguage := "Pastikan semua informasi nama negara dan kota disajikan dalam Bahasa Indonesia."
	prompt := fmt.Sprintf(`Cari informasi tentang destinasi "%s". %s`, query, instructionLanguage)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gagal generate content dari Gemini: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("respon dari Gemini kosong atau tidak sesuai format yang diharapkan")
	}

	jsonTextPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		if blob, isBlob := resp.Candidates[0].Content.Parts[0].(genai.Blob); isBlob && blob.MIMEType == "application/json" {
			jsonTextPart = genai.Text(blob.Data)
		} else {
			return nil, fmt.Errorf("format part respon tidak terduga dari Gemini: %T. Diharapkan genai.Text", resp.Candidates[0].Content.Parts[0])
		}
	}

	rawJSON := cleanJSONResponse(string(jsonTextPart))

	var searchResult struct {
		Results []PlaceSearch `json:"results"`
	}

	if err := json.Unmarshal([]byte(rawJSON), &searchResult); err != nil {
		return nil, fmt.Errorf("gagal parsing JSON dari teks respon Gemini: %w. \nJSON Mentah: %s", err, rawJSON)
	}

	return searchResult.Results, nil
}

// Fitur GetPlaces
func (s GeminiService) GetPlaces(preferences []string, country, town string) (map[string]interface{}, error) {
	apiKey := config.GeminiAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("API Key tidak ditemukan")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")

	responseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"accomodations": {
				Type:        genai.TypeArray,
				Description: "List of recommended accommodations.",
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"name": {Type: genai.TypeString, Description: "Name of the accommodation."},
						"town": {Type: genai.TypeString, Description: "Town where the accommodation is located."},
					},
					Required: []string{"name", "town"},
				},
			},
			"places": {
				Type:        genai.TypeArray,
				Description: "List of recommended tourist places.",
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"description": {Type: genai.TypeString, Description: "Brief description of the place."},
						"name":        {Type: genai.TypeString, Description: "Name of the place."},
						"town":        {Type: genai.TypeString, Description: "Town where the place is located."},
						"type":        {Type: genai.TypeString, Description: "Type of place (e.g., Food, Historical Site, Cultural Site)."},
					},
					Required: []string{"description", "name", "town", "type"},
				},
			},
		},
		Required: []string{"accomodations", "places"},
	}

	model.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	}

	instructionLanguage := "Pastikan semua informasi yang Anda berikan, termasuk nama tempat, nama akomodasi, deskripsi, dan tipe/kategori, disajikan dalam Bahasa Indonesia."

	promptPreferences := fmt.Sprintf("%v", preferences)
	prompt := fmt.Sprintf("Berikan daftar tempat wisata dan akomodasi di %s, %s berdasarkan preferensi berikut: %s.%s",
		town, country, promptPreferences, instructionLanguage)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gagal generate content dari Gemini: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("respon dari Gemini kosong atau tidak sesuai format yang diharapkan")
	}

	jsonTextPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		if blob, isBlob := resp.Candidates[0].Content.Parts[0].(genai.Blob); isBlob && blob.MIMEType == "application/json" {
			jsonTextPart = genai.Text(blob.Data)
		} else {
			return nil, fmt.Errorf("format part respon tidak terduga dari Gemini: %T. Diharapkan genai.Text", resp.Candidates[0].Content.Parts[0])
		}
	}

	rawJSON := removeMarkdownCodeBlock(string(jsonTextPart))

	var geminiResult GeminiResponseGetPlaces
	if err := json.Unmarshal([]byte(rawJSON), &geminiResult); err != nil {
		return nil, fmt.Errorf("gagal parsing JSON dari teks respon Gemini: %w. \nJSON Mentah: %s", err, rawJSON)
	}

	for i := range geminiResult.Places {
		name := geminiResult.Places[i].Name
		imageURLs, err := GetGoogleImagesPlaces(name)
		if err == nil {
			geminiResult.Places[i].Image = imageURLs
		}
	}

	for i := range geminiResult.Accomodations {
		name := geminiResult.Accomodations[i].Name
		imageURLs, err := GetGoogleImages(name)
		if err == nil {
			geminiResult.Accomodations[i].Image = imageURLs
		}
	}

	response := map[string]interface{}{
		"accomodations": geminiResult.Accomodations,
		"places":        geminiResult.Places,
	}

	return response, nil
}

// fitur timeline
func (s GeminiService) GetTimeline(accommodation, town, country, startDate, endDate string,
	places []SelectedPlace,
) (*TimelineResponse, error) {
	apiKey := config.GeminiAPIKey
	if apiKey == "" {
		return nil, errors.New("API Key tidak ditemukan")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")

	placeVisitedDetailSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"type":     {Type: genai.TypeString, Description: "Jenis tempat (misalnya: Hotel, Restoran, Tempat Wisata Sejarah, Taman, dll.)."},
			"landmark": {Type: genai.TypeString, Description: "Nama spesifik tempat atau landmark."},
			"roadName": {Type: genai.TypeString, Description: "Nama jalan lokasi. Harus diisi dan relevan."},
			"town":     {Type: genai.TypeString, Description: "Nama kota atau daerah tempat ini berada."},
			"time":     {Type: genai.TypeString, Description: "Waktu kunjungan atau kegiatan dalam format HH:MM (misalnya: 14:00)."},
		},
		Required: []string{"type", "landmark", "roadName", "town", "time"},
	}

	dailyVisitSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"date": {Type: genai.TypeString, Description: "Tanggal kunjungan dalam format YYYY-MM-DD."},
			"data": {
				Type:        genai.TypeArray,
				Description: "Daftar detail tempat yang dikunjungi atau kegiatan pada tanggal tersebut.",
				Items:       placeVisitedDetailSchema,
			},
		},
		Required: []string{"date", "data"},
	}

	timelineResponseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"budget":    {Type: genai.TypeString, Description: "Estimasi total budget perjalanan dalam IDR (misalnya: 1000000)."},
			"town":      {Type: genai.TypeString, Description: "Nama kota utama untuk rencana perjalanan ini."},
			"country":   {Type: genai.TypeString, Description: "Nama negara untuk rencana perjalanan ini."},
			"startDate": {Type: genai.TypeString, Description: "Tanggal mulai perjalanan dalam format YYYY-MM-DD."},
			"endDate":   {Type: genai.TypeString, Description: "Tanggal akhir perjalanan dalam format YYYY-MM-DD."},
			"placeVisited": {
				Type:        genai.TypeArray,
				Description: "Rincian rencana perjalanan per hari.",
				Items:       dailyVisitSchema,
			},
		},
		Required: []string{"budget", "town", "country", "startDate", "endDate", "placeVisited"},
	}

	model.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   timelineResponseSchema,
	}

	var placesStrBuilder strings.Builder
	if len(places) > 0 {
		placesStrBuilder.WriteString("Berikut adalah preferensi waktu dan daftar tempat yang ingin dikunjungi: ")
		for i, sp := range places {
			if i > 0 {
				placesStrBuilder.WriteString("; ")
			}
			placesStrBuilder.WriteString(fmt.Sprintf("Waktu '%s', tempat: [%s]", sp.TimeOfDay, strings.Join(sp.Places, ", ")))
		}
		placesStrBuilder.WriteString(".") // Akhiri kalimat
	} else {
		placesStrBuilder.WriteString("Tidak ada preferensi tempat spesifik yang diberikan untuk dipertimbangkan secara khusus.")
	}

	instructionLanguage := "Semua informasi, termasuk nama tempat, tipe, deskripsi jalan, kota, dan estimasi budget, harus dalam Bahasa Indonesia."
	prompt := fmt.Sprintf(
		"Buatkan rencana perjalanan dari akomodasi: %s, di kota/daerah: %s, negara: %s, dari tanggal %s hingga %s. "+
			"%s "+
			"Sertakan juga estimasi budget keseluruhan dalam IDR. %s",
		accommodation, town, country, startDate, endDate,
		placesStrBuilder.String(),
		instructionLanguage,
	)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gagal generate content dari Gemini: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("respon dari Gemini kosong atau tidak sesuai format yang diharapkan")
	}

	jsonTextPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		if blob, isBlob := resp.Candidates[0].Content.Parts[0].(genai.Blob); isBlob && blob.MIMEType == "application/json" {
			jsonTextPart = genai.Text(blob.Data)
		} else {
			return nil, fmt.Errorf("format part respon tidak terduga dari Gemini: %T. Diharapkan genai.Text", resp.Candidates[0].Content.Parts[0])
		}
	}

	rawJSON := removeMarkdownCodeBlock(string(jsonTextPart))

	var result TimelineResponse
	if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
		return nil, fmt.Errorf("gagal parsing JSON dari teks respon Gemini: %w. \nJSON Mentah: %s", err, rawJSON)
	}

	for di, daily := range result.PlaceVisited {
		for pi, place := range daily.Data {
			name := place.Landmark
			imageURLs, err := GetGoogleImagesPlaces(name)
			if err == nil && len(imageURLs) >= 2 {
				result.PlaceVisited[di].Data[pi].Image = imageURLs[:2]
			}
		}
	}

	return &result, nil

}

// Fitur GetPlaceDetail
func (s GeminiService) GetPlaceDetail(placeType, landmark, town string) (map[string]interface{}, error) {
	apiKey := config.GeminiAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("API Key tidak ditemukan")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")

	responseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"description": {
				Type:        genai.TypeString,
				Description: "Deskripsi detail dan informatif tentang tempat atau landmark yang diminta, dalam Bahasa Indonesia.",
			},
		},
		Required: []string{"description"},
	}

	model.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	}
	instructionLanguage := "Semua informasi harus dalam Bahasa Indonesia."
	prompt := fmt.Sprintf(`Berikan deskripsi detail tentang %s bernama "%s" yang terletak di %s. %s`,
		placeType, landmark, town, instructionLanguage)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gagal generate content dari Gemini: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("respon dari Gemini kosong atau tidak sesuai format yang diharapkan")
	}

	jsonTextPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		if blob, isBlob := resp.Candidates[0].Content.Parts[0].(genai.Blob); isBlob && blob.MIMEType == "application/json" {
			jsonTextPart = genai.Text(blob.Data)
		} else {
			return nil, fmt.Errorf("format part respon tidak terduga dari Gemini: %T. Diharapkan genai.Text", resp.Candidates[0].Content.Parts[0])
		}
	}

	rawJSON := removeMarkdownCodeBlock(string(jsonTextPart))

	var aiResponse struct {
		Description string `json:"description"`
	}
	if err := json.Unmarshal([]byte(rawJSON), &aiResponse); err != nil {
		return nil, fmt.Errorf("gagal parsing JSON dari teks respon Gemini: %w. \nJSON Mentah: %s", err, rawJSON)
	}

	var detail LandmarkDetail
	detail.Description = aiResponse.Description

	images, imgErr := GetGoogleImagesPlaces(landmark)
	if imgErr == nil && len(images) > 0 {
		if len(images) >= 2 {
			detail.Images = images[:2]
		} else {
			detail.Images = images
		}
	}

	return map[string]interface{}{
		"description": detail.Description,
		"images":      detail.Images,
	}, nil
}
