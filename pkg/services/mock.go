package services

import "encoding/json"

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
	raw := []byte(`{
        "budget": "5000000",
        "town": "Paris",
        "country": "France",
        "startDate": "2024-08-12",
        "endDate": "2024-08-16",
        "placeVisited": [
            {
                "date": "2024-08-12",
                "data": [
                    {
                        "type": "Hotel",
                        "landmark": "Hotel Paris",
                        "roadName": "To be determined",
                        "town": "Paris",
                        "time": "14:00",
                        "image": [
                            "https://media-cdn.tripadvisor.com/media/photo-s/2a/01/a6/b0/facade-solly-hotel-paris.jpg",
                            "https://anoushkaprobyn.com/wp-content/uploads/2024/02/Disneyland-Hotel-Paris-Review2.jpg"
                        ]
                    }
                ]
            },
            {
                "date": "2024-08-13",
                "data": [
                    {
                        "type": "Landmark",
                        "landmark": "Menara Eiffel",
                        "roadName": "Champ de Mars, 5 Avenue Anatole France",
                        "town": "Paris",
                        "time": "09:00",
                        "image": [
                            "https://i.pinimg.com/474x/78/c4/58/78c45898090779bb3bf0b37b7ac2bcfe.jpg",
                            "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Tour_Eiffel_Wikimedia_Commons.jpg/250px-Tour_Eiffel_Wikimedia_Commons.jpg"
                        ]
                    },
                    {
                        "type": "Museum",
                        "landmark": "Centre Georges-Pompidou",
                        "roadName": "Place Georges-Pompidou",
                        "town": "Paris",
                        "time": "14:00",
                        "image": [
                            "https://images.adsttc.com/media/images/6515/ba44/7316/3253/1be0/7822/newsletter/architecture-classics-centre-georges-pompidou-renzo-piano-building-workshop-plus-richard-rogers_11.jpg?1695922762",
                            "https://www.centrepompidou.fr/fileadmin/_processed_/2/4/csm_collection-notrebatiment-photofacaderuerambuteau2021_1920x750_af5ca8a213.jpg"
                        ]
                    },
                    {
                        "type": "Museum",
                        "landmark": "Musée Picasso",
                        "roadName": "5 Rue de Thorigny",
                        "town": "Paris",
                        "time": "19:00",
                        "image": [
                            "https://www.museepicassoparis.fr/sites/default/files/2020-01/Horaires---Musee-Picasso---Voyez-Vous-%28Chloe-Vollmer-Lo%29--8697.jpg",
                            "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/H%C3%B4tel_Sal%C3%A9.JPG/1200px-H%C3%B4tel_Sal%C3%A9.JPG"
                        ]
                    }
                ]
            },
            {
                "date": "2024-08-14",
                "data": [
                    {
                        "type": "Landmark",
                        "landmark": "Menara Eiffel",
                        "roadName": "Champ de Mars, 5 Avenue Anatole France",
                        "town": "Paris",
                        "time": "09:00",
                        "image": [
                            "https://i.pinimg.com/474x/78/c4/58/78c45898090779bb3bf0b37b7ac2bcfe.jpg",
                            "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Tour_Eiffel_Wikimedia_Commons.jpg/250px-Tour_Eiffel_Wikimedia_Commons.jpg"
                        ]
                    },
                    {
                        "type": "Museum",
                        "landmark": "Centre Georges-Pompidou",
                        "roadName": "Place Georges-Pompidou",
                        "town": "Paris",
                        "time": "14:00",
                        "image": [
                            "https://images.adsttc.com/media/images/6515/ba44/7316/3253/1be0/7822/newsletter/architecture-classics-centre-georges-pompidou-renzo-piano-building-workshop-plus-richard-rogers_11.jpg?1695922762",
                            "https://www.centrepompidou.fr/fileadmin/_processed_/2/4/csm_collection-notrebatiment-photofacaderuerambuteau2021_1920x750_af5ca8a213.jpg"
                        ]
                    },
                    {
                        "type": "Park",
                        "landmark": "Parc de la Villette",
                        "roadName": "211 Avenue Jean Jaurès",
                        "town": "Paris",
                        "time": "19:00",
                        "image": [
                            "https://upload.wikimedia.org/wikipedia/commons/e/ed/Rio_Samba_School_statue_%40_Parc_de_La_Villette_%40_Paris_%2828881779791%29.jpg",
                            "https://www.tschumi.com/cms/assets/2c88af60-b8d2-4044-92c7-8806472ecc3d?width=700&fit=contain"
                        ]
                    }
                ]
            },
            {
                "date": "2024-08-15",
                "data": [
                    {
                        "type": "Landmark",
                        "landmark": "Menara Eiffel",
                        "roadName": "Champ de Mars, 5 Avenue Anatole France",
                        "town": "Paris",
                        "time": "09:00",
                        "image": [
                            "https://i.pinimg.com/474x/78/c4/58/78c45898090779bb3bf0b37b7ac2bcfe.jpg",
                            "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Tour_Eiffel_Wikimedia_Commons.jpg/250px-Tour_Eiffel_Wikimedia_Commons.jpg"
                        ]
                    },
                    {
                        "type": "Museum",
                        "landmark": "Centre Georges-Pompidou",
                        "roadName": "Place Georges-Pompidou",
                        "town": "Paris",
                        "time": "14:00",
                        "image": [
                            "https://images.adsttc.com/media/images/6515/ba44/7316/3253/1be0/7822/newsletter/architecture-classics-centre-georges-pompidou-renzo-piano-building-workshop-plus-richard-rogers_11.jpg?1695922762",
                            "https://www.centrepompidou.fr/fileadmin/_processed_/2/4/csm_collection-notrebatiment-photofacaderuerambuteau2021_1920x750_af5ca8a213.jpg"
                        ]
                    },
                    {
                        "type": "Museum",
                        "landmark": "Musée Picasso",
                        "roadName": "5 Rue de Thorigny",
                        "town": "Paris",
                        "time": "19:00",
                        "image": [
                            "https://www.museepicassoparis.fr/sites/default/files/2020-01/Horaires---Musee-Picasso---Voyez-Vous-%28Chloe-Vollmer-Lo%29--8697.jpg",
                            "https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/H%C3%B4tel_Sal%C3%A9.JPG/1200px-H%C3%B4tel_Sal%C3%A9.JPG"
                        ]
                    }
                ]
            },
            {
                "date": "2024-08-16",
                "data": [
                    {
                        "type": "Hotel",
                        "landmark": "Hotel Paris",
                        "roadName": "To be determined",
                        "town": "Paris",
                        "time": "11:00",
                        "image": [
                            "https://media-cdn.tripadvisor.com/media/photo-s/2a/01/a6/b0/facade-solly-hotel-paris.jpg",
                            "https://anoushkaprobyn.com/wp-content/uploads/2024/02/Disneyland-Hotel-Paris-Review2.jpg"
                        ]
                    }
                ]
            }
        ]
    }`)

	var resp TimelineResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return &resp, nil

}
