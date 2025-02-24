package database

import (
	"fmt"
	"healcationBackend/models"
	"time"
)

func SeedHistory() {
	var count int64
	DB.Model(&models.History{}).Count(&count)

	if count > 0 {
		fmt.Println("History already seeded, skipping...")
		return
	}

	history := models.History{
		Country:     "Indonesia",
		Town:        "Bali",
		StartDate:   time.Date(2023, time.January, 1, 7, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2023, time.January, 7, 7, 0, 0, 0, time.UTC),
		Image:       "https://example.com/bali.jpg",
		Description: "Bali Vacation",
		SelectedAccomodation: []models.SelectedAccomodation{
			{Name: "Hotel Entah", Image: "https://example.com/hotel.jpg"},
		},
		SelectedPlaces: []models.SelectedPlace{
			{PlaceToVisit: "Temple", Town: "Somewhere in Bali", Image: "https://example.com/temple.jpg"},
		},
	}

	if err := DB.Create(&history).Error; err != nil {
		fmt.Println("Error seeding history:", err)
	} else {
		fmt.Println("History seeded successfully")
	}
}

// func SeedUser() {

// var count int64
// DB.Model(&models.User{}).Count(&count)

// if count > 0 {
// 	fmt.Println("User already seeded, skipping...")
// 	return
// }
// 	user := models.User{
// 		Username: "admin",
// 		Email:    "admin@gmai.com",
// 		Password: "admin123",
// 	}

// 	if err := DB.Create(&user).Error; err != nil {
// 		fmt.Println("Error seeding user:", err)
// 	} else {
// 		fmt.Println("User seeded successfully")
// 	}
// }

func Seed() {
	SeedHistory()
	// SeedUser()
}
