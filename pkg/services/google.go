package services

import (
	"encoding/json"
	"fmt"
	"healcationBackend/pkg/config"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func validateImageURL(client *http.Client, imageURL string) bool {
	if !(strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://")) {
		return false
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := client.Head(imageURL)
	if err == nil {
		defer res.Body.Close()
		if res.StatusCode == http.StatusOK {
			if strings.HasPrefix(res.Header.Get("Content-Type"), "image/") {
				return true
			}
		}
	}
	resp, err := client.Get(imageURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close()
	buf := make([]byte, 512)
	n, _ := io.ReadFull(resp.Body, buf)
	return strings.HasPrefix(http.DetectContentType(buf[:n]), "image/")
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

	client := &http.Client{
		Timeout: 8 * time.Second,
	}
	validImageURLs := []string{}
	const (
		targetCount = 1
		maxAttempts = 5
		baseFetch   = 2
	)

	for attempt := 0; attempt < maxAttempts && len(validImageURLs) < targetCount; attempt++ {
		numToFetch := baseFetch + attempt
		if numToFetch > 5 {
			numToFetch = 5
		}

		items, err := fetchRawGoogleImages(query, apiKey, cx, numToFetch, 1)
		if err != nil || len(items) < 1 {
			break
		}

		links := make([]string, len(items))
		for i, itm := range items {
			links[i] = itm.Link
		}
		fmt.Printf("[Attempt %d] fetched links: %v\n", attempt+1, links)

		start := 0
		if len(items) > 1 {
			start = len(items) - 1
		}
		tailItems := items[start:]

		for _, item := range tailItems {
			if validateImageURL(client, item.Link) {
				dup := false
				for _, u := range validImageURLs {
					if u == item.Link {
						dup = true
						break
					}
				}
				if !dup {
					validImageURLs = append(validImageURLs, item.Link)
					if len(validImageURLs) == targetCount {
						break
					}
				}
			}
		}
	}

	if len(validImageURLs) < targetCount {
		return nil, fmt.Errorf("gagal menemukan %d gambar valid untuk '%s' setelah %d percobaan",
			targetCount, query, maxAttempts)
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

	client := &http.Client{
		Timeout: 8 * time.Second,
	}
	validImageURLs := []string{}
	const (
		targetCount = 2
		maxAttempts = 5
		baseFetch   = 3
	)

	for attempt := 0; attempt < maxAttempts && len(validImageURLs) < targetCount; attempt++ {
		numToFetch := baseFetch + attempt
		if numToFetch > 5 {
			numToFetch = 5
		}

		items, err := fetchRawGoogleImages(query, apiKey, cx, numToFetch, 1)
		if err != nil || len(items) < 1 {
			break
		}

		links := make([]string, len(items))
		for i, itm := range items {
			links[i] = itm.Link
		}
		fmt.Printf("[Attempt %d] fetched links: %v\n", attempt+1, links)

		start := 0
		if len(items) > 2 {
			start = len(items) - 2
		}
		tailItems := items[start:]

		for _, item := range tailItems {
			if validateImageURL(client, item.Link) {
				dup := false
				for _, u := range validImageURLs {
					if u == item.Link {
						dup = true
						break
					}
				}
				if !dup {
					validImageURLs = append(validImageURLs, item.Link)
					if len(validImageURLs) == targetCount {
						break
					}
				}
			}
		}
	}

	if len(validImageURLs) < targetCount {
		return nil, fmt.Errorf("gagal menemukan %d gambar valid untuk '%s' setelah %d percobaan",
			targetCount, query, maxAttempts)
	}
	return validImageURLs, nil
}
