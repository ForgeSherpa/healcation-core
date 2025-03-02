package database

import (
	"log"
	"time"

	"healcationBackend/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedHistory() {
	var count int64
	DB.Model(&models.History{}).Count(&count)

	if count > 0 {
		log.Println("History already seeded, skipping...")
		return
	}

	history := models.History{
		Country:     "Indonesia",
		Town:        "Bali",
		StartDate:   time.Date(2023, time.January, 1, 7, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2023, time.January, 7, 7, 0, 0, 0, time.UTC),
		Image:       models.StringArray{"https://example.com/bali.jpg"},
		Description: "Bali Vacation",
		SelectedAccomodation: []models.SelectedAccomodation{
			{Name: "Hotel Entah", Image: models.StringArray{"https://example.com/hotel.jpg"}},
		},
		SelectedPlaces: []models.SelectedPlace{
			{PlaceToVisit: "Temple", Town: "Somewhere in Bali", Image: models.StringArray{"https://example.com/temple.jpg"}},
		},
	}

	if err := DB.Create(&history).Error; err != nil {
		log.Println("Error seeding history:", err)
	} else {
		log.Println("History seeded successfully")
	}
}

func SeedUser() {
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count > 0 {
		log.Println("User already seeded, skipping...")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("delvin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	user := models.User{
		Username: "delvin",
		Email:    "delvin@gmail.com",
		Password: string(hashedPassword),
	}

	if err := DB.Create(&user).Error; err != nil {
		log.Println("Error seeding user:", err)
	} else {
		log.Println("User seeded successfully")
	}
}

func Seed() {
	SeedHistory()
	SeedUser()
}
