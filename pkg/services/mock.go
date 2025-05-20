package services

import "fmt"

type GeminiMockService struct{}

func NewGeminiMockService() AIService {
	return GeminiMockService{}
}

func (m GeminiMockService) Search(query string) ([]PlaceSearch, error) {
	results := []PlaceSearch{
		{Country: "Indonesia", Town: "Denpasar"},
		{Country: "Indonesia", Town: "Ubud"},
		{Country: "Indonesia", Town: "Seminyak"},
		{Country: "Indonesia", Town: "Kuta"},
		{Country: "Indonesia", Town: "Nusa Dua"},
	}
	return results, nil
}

func (m GeminiMockService) GetPlaces(preferences []string, country, town string) (map[string]interface{}, error) {
	accommodations := []map[string]interface{}{
		{
			"image": []string{
				"https://ecom-cvweb.s3-us-west-2.amazonaws.com/s3fs-public/styles/hotel_hero_image/public/le-bristol-paris-v46-3-16-2017.jpg?itok=xPUyqeTe",
			},
			"name": "Le Bristol Paris - Oetker Collection",
		},
		{
			"image": []string{
				"https://www.dorchestercollection.com/media/ipgftz4o/h%C3%B4tel-plaza-ath%C3%A9n%C3%A9e-la-cour-jardin-courtyard-garden-dorchester-collection.jpg?rxy=0.47118440639507575%2C0.8476572327174599&width=746&height=810&rmode=crop",
			},
			"name": "Hotel Plaza Athénée - Dorchester Collection",
		},
	}

	places := []map[string]interface{}{
		{
			"description": "Museum seni terbesar dan paling terkenal di dunia, rumah bagi Mona Lisa dan karya seni klasik lainnya.",
			"image": []string{
				"https://api-www.louvre.fr/sites/default/files/2021-01/cour-napoleon-et-pyramide_1.jpg",
				"https://upload.wikimedia.org/wikipedia/commons/thumb/6/66/Louvre_Museum_Wikimedia_Commons.jpg/800px-Louvre_Museum_Wikimedia_Commons.jpg",
			},
			"name": "Musée du Louvre",
			"town": "paris",
			"type": "Museum",
		},
		{
			"description": "Ikon kota Paris, menawarkan pemandangan panorama kota yang menakjubkan.",
			"image": []string{
				"https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg/640px-Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg",
				"https://i.natgeofe.com/k/c41b4f59-181c-4747-ad20-ef69987c8d59/eiffel-tower-night.jpg?wp=1&w=1084.125&h=1627.5",
			},
			"name": "Eiffel Tower",
			"town": "paris",
			"type": "Landmark",
		},
	}

	result := map[string]interface{}{
		"accomodations": accommodations,
		"places":        places,
	}
	return result, nil
}

func (m GeminiMockService) GetPlaceDetail(placeType, landmark, town string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"description": "Menara Eiffel adalah menara besi kisi yang terletak di Champ de Mars, Paris, Prancis. Dibangun pada tahun 1889 sebagai pintu masuk untuk Eksposisi Universelle, ia telah menjadi ikon global Prancis dan salah satu struktur paling dikenal di dunia. Pengunjung dapat naik ke atas untuk menikmati pemandangan kota Paris yang spektakuler.",
		"images": []string{
			"https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg/1200px-Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg",
			"https://www.toureiffel.paris/themes/custom/tour_eiffel/build/images/home-discover-bg.jpg",
		},
	}, nil
}

func (m GeminiMockService) GetTimeline(accommodation, town, country, startDate, endDate string,
	places []SelectedPlace,
) (*TimelineResponse, error) {
	fmt.Print("Mocking GetTimeline")
	return nil, nil
}
