package services

import (
	"encoding/json"
	"fmt"
	"healcationBackend/pkg/config"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func isValidURL(link string) bool {
	return strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")
}

func fetchRawGoogleImages(query, apiKey, cx string, numToFetch, startIndex int) ([]struct {
	Link string `json:"link"`
}, error) {
	if numToFetch <= 0 {
		return []struct {
			Link string `json:"link"`
		}{}, nil
	}
	apiURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&key=%s&cx=%s&searchType=image&num=%d&start=%d",
		url.QueryEscape(query), apiKey, cx, numToFetch, startIndex)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan permintaan ke Google API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gagal mendapatkan gambar dari Google: status %d, body: %s", resp.StatusCode, string(errorBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca body respons: %w", err)
	}

	var googleResp struct {
		Items []struct {
			Link string `json:"link"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &googleResp); err != nil {
		return nil, fmt.Errorf("gagal unmarshal JSON respons: %w. Respons body: %s", err, string(body))
	}
	return googleResp.Items, nil
}

// return 1 image URL
func GetGoogleImages(query string) ([]string, error) {
	apiKey := config.GoogleAPIKey
	cx := config.GoogleAPI_CX

	if apiKey == "" || cx == "" {
		return nil, fmt.Errorf("google API Key atau CX tidak ditemukan")
	}

	validImageURLs := []string{}
	const targetCount = 1

	startIndexForAPICall := 1

	maxAttempts := 5
	attempts := 0

	for len(validImageURLs) < targetCount && attempts < maxAttempts {
		attempts++

		neededValidURLs := targetCount - len(validImageURLs)

		numToFetchThisAttempt := neededValidURLs
		if attempts == 1 {
			numToFetchThisAttempt = targetCount + 1
		} else {
			numToFetchThisAttempt = neededValidURLs + 2
		}
		if numToFetchThisAttempt > 5 {
			numToFetchThisAttempt = 5
		}

		items, err := fetchRawGoogleImages(query, apiKey, cx, numToFetchThisAttempt, startIndexForAPICall)
		if err != nil {
			break
		}

		if len(items) == 0 {
			break
		}

		itemsProcessedFromThisFetch := 0
		for _, item := range items {
			itemsProcessedFromThisFetch++
			if isValidURL(item.Link) {
				isDuplicate := false
				for _, existingURL := range validImageURLs {
					if existingURL == item.Link {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					validImageURLs = append(validImageURLs, item.Link)
					if len(validImageURLs) == targetCount {
						break
					}
				}
			}
		}

		startIndexForAPICall += itemsProcessedFromThisFetch

		if len(validImageURLs) == targetCount {
			break
		}
	}

	if len(validImageURLs) == 0 {
		return nil, fmt.Errorf("tidak ada URL gambar HTTP/HTTPS yang valid ditemukan untuk query: '%s' setelah %d percobaan", query, attempts)
	}

	return validImageURLs, nil
}

// return 2 image URL
func GetGoogleImagesPlaces(query string) ([]string, error) {
	apiKey := config.GoogleAPIKey
	cx := config.GoogleAPI_CX

	if apiKey == "" || cx == "" {
		return nil, fmt.Errorf("google API Key atau CX tidak ditemukan")
	}

	validImageURLs := []string{}
	const targetCount = 2

	startIndexForAPICall := 1

	maxAttempts := 5
	attempts := 0

	for len(validImageURLs) < targetCount && attempts < maxAttempts {
		attempts++

		neededValidURLs := targetCount - len(validImageURLs)

		numToFetchThisAttempt := neededValidURLs
		if attempts == 1 {
			numToFetchThisAttempt = targetCount + 1
		} else {
			numToFetchThisAttempt = neededValidURLs + 2
		}
		if numToFetchThisAttempt > 5 {
			numToFetchThisAttempt = 5
		}

		items, err := fetchRawGoogleImages(query, apiKey, cx, numToFetchThisAttempt, startIndexForAPICall)
		if err != nil {
			break
		}

		if len(items) == 0 {
			break
		}

		itemsProcessedFromThisFetch := 0
		for _, item := range items {
			itemsProcessedFromThisFetch++
			if isValidURL(item.Link) {
				isDuplicate := false
				for _, existingURL := range validImageURLs {
					if existingURL == item.Link {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					validImageURLs = append(validImageURLs, item.Link)
					if len(validImageURLs) == targetCount {
						break
					}
				}
			}
		}

		startIndexForAPICall += itemsProcessedFromThisFetch

		if len(validImageURLs) == targetCount {
			break
		}
	}

	if len(validImageURLs) == 0 {
		return nil, fmt.Errorf("tidak ada URL gambar HTTP/HTTPS yang valid ditemukan untuk query: '%s' setelah %d percobaan", query, attempts)
	}

	return validImageURLs, nil
}
