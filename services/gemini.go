package services

import "fmt"

func SearchGemini(query string) []map[string]string {
	staticResults := []map[string]string{
		{
			"town":    "Bali",
			"country": "Indonesia",
		},
	}

	fmt.Println("Search query:", query)

	return staticResults
}
