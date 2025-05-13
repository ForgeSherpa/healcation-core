package services

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
) (map[string]interface{}, error) {
	timeline := map[string][]map[string]string{
		"2024-08-12": {
			{
				"image":    "https://www.deluxefrance.com/public/img/big/ParisCDGairportwithplanejpg_652951e7cfa0c4.73613563.jpg",
				"landmark": "Arrival in Paris",
				"roadName": "Charles de Gaulle Airport (CDG)",
				"time":     "Afternoon",
				"town":     "Paris",
				"type":     "Arrival",
			},
			{
				"image":    "https://www.cataloniahotels.com/en/blog/wp-content/uploads/2022/03/check-in-hotel.jpg",
				"landmark": "Check in to Hotel",
				"roadName": "Rue de Rivoli",
				"time":     "Evening",
				"town":     "Paris",
				"type":     "Accommodation",
			},
			{
				"image":    "https://media.tacdn.com/media/attractions-splice-spp-674x446/07/8f/2d/63.jpg",
				"landmark": "Dinner near the Seine",
				"roadName": "Quai des Grands Augustins",
				"time":     "Night",
				"town":     "Paris",
				"type":     "Restaurant",
			},
		},
		"2024-08-13": {
			{
				"image":    "https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg/640px-Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg",
				"landmark": "Eiffel Tower",
				"roadName": "Champ de Mars",
				"time":     "Morning",
				"town":     "Paris",
				"type":     "Landmark",
			},
			{
				"image":    "https://upload.wikimedia.org/wikipedia/commons/thumb/2/29/MG-Paris-Champ_de_Mars.jpg/1200px-MG-Paris-Champ_de_Mars.jpg",
				"landmark": "Champ de Mars",
				"roadName": "Avenue Anatole France",
				"time":     "Afternoon",
				"town":     "Paris",
				"type":     "Park",
			},
			{
				"image":    "https://www.cruiseontheseine.com/wp-content/uploads/2023/11/seine-river-cruise-1024x686.jpg",
				"landmark": "Seine River Cruise",
				"roadName": "Port de la Bourdonnais",
				"time":     "Late Afternoon",
				"town":     "Paris",
				"type":     "Activity",
			},
			{
				"image":    "https://dynamic-media-cdn.tripadvisor.com/media/photo-o/03/51/9c/c2/chartier.jpg?w=900&h=500&s=1",
				"landmark": "Dinner at Le Bouillon Chartier",
				"roadName": "Rue du Faubourg Montmartre",
				"time":     "Night",
				"town":     "Paris",
				"type":     "Restaurant",
			},
		},
	}

	result := map[string]interface{}{
		"budget":   "Rp 7.500.000 - Rp 15.000.000",
		"country":  "France Static",
		"timeline": timeline,
		"title":    "Parisian Adventure: A 5-Day Itinerary",
		"town":     "Paris",
	}
	return result, nil
}
