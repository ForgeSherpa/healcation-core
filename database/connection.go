package database

import (
	"log"
	"os"

	"healcationBackend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	dbPath := os.Getenv("DB_PATH")
	log.Println("Connecting to SQLite database:", dbPath)

	connection, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	DB = connection

	log.Println("Running database migrations...")
	if err := connection.AutoMigrate(
		&models.History{},
		&models.User{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migrations completed")

	Seed()
	log.Println("Seed data inserted successfully")
}
