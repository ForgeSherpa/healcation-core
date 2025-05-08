package services

import (
	"encoding/json"
	"fmt"
	"healcationBackend/pkg/config"
	"io"
	"net/http"
	"net/url"
)

func GetGoogleImages(query string) ([]string, error) {
	apiKey := config.GoogleAPIKey
	cx := config.GoogleAPI_CX

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
	apiKey := config.GoogleAPIKey
	cx := config.GoogleAPI_CX

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
