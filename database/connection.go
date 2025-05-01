package database

import (
	"fmt"
	"log"
	"os"

	"healcationBackend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var DB *gorm.DB

func Connect() {
	primaryURL := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	if primaryURL == "" || authToken == "" {
		log.Fatal("TURSO_DATABASE_URL dan TURSO_AUTH_TOKEN harus di-set")
	}

	dsn := fmt.Sprintf("%s?authToken=%s", primaryURL, authToken)

	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "libsql",
		DSN:        dsn,
	}, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to Turso database: %v", err)
	}
	log.Println("Database connected successfully")

	DB = db

	if err := DB.AutoMigrate(&models.History{}, &models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrations completed")

	// Seed()
	// log.Println("Seed data inserted successfully")
}
