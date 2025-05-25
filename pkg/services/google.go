package services

import (
	"encoding/json"
	"fmt"
	"healcationBackend/pkg/config"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

var (
	headCount int32
	getCount  int32
)

func validateImageURL(client *http.Client, imageURL string) bool {
	if !(strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://")) {
		return false
	}

	atomic.AddInt32(&headCount, 1)
	fmt.Printf("[validateImageURL] HEAD call #%d to %s\n", atomic.LoadInt32(&headCount), imageURL)

	headRes, err := client.Head(imageURL)
	if err == nil {
		defer headRes.Body.Close()
		if headRes.StatusCode == http.StatusOK {
			ct := headRes.Header.Get("Content-Type")
			if strings.HasPrefix(ct, "image/") {
				return true
			}
		}
	}

	atomic.AddInt32(&getCount, 1)
	fmt.Printf("[validateImageURL] GET  call #%d to %s\n", atomic.LoadInt32(&getCount), imageURL)
	getRes, err := client.Get(imageURL)
	if err != nil {
		return false
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		return false
	}

	buf := make([]byte, 512)
	n, err := io.ReadFull(getRes.Body, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return false
	}
	detected := http.DetectContentType(buf[:n])
	return strings.HasPrefix(detected, "image/")
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
	const targetCount = 1

	startIndex := 1

	maxAttempts := 5
	for attempts := 0; attempts < maxAttempts && len(validImageURLs) < targetCount; attempts++ {
		needed := targetCount - len(validImageURLs)
		perCall := needed + 2
		if perCall > 5 {
			perCall = 5
		}

		items, err := fetchRawGoogleImages(query, apiKey, cx, perCall, startIndex)
		if err != nil || len(items) == 0 {
			break
		}

		processed := 0
		for _, item := range items {
			processed++
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

		startIndex += processed
	}

	if len(validImageURLs) == 0 {
		return nil, fmt.Errorf("tidak ada URL gambar valid untuk query '%s' setelah %d percobaan", query, maxAttempts)
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
	const targetCount = 2

	startIndex := 1
	maxAttempts := 5

	for attempts := 0; attempts < maxAttempts && len(validImageURLs) < targetCount; attempts++ {
		needed := targetCount - len(validImageURLs)
		perCall := needed + 2
		if perCall > 5 {
			perCall = 5
		}

		items, err := fetchRawGoogleImages(query, apiKey, cx, perCall, startIndex)
		if err != nil || len(items) == 0 {
			break
		}

		processed := 0
		for _, item := range items {
			processed++
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

		startIndex += processed
	}

	if len(validImageURLs) == 0 {
		return nil, fmt.Errorf("tidak ada URL gambar valid untuk query '%s' setelah %d percobaan", query, maxAttempts)
	}
	return validImageURLs, nil
}
