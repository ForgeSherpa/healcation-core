package database

import (
	"fmt"
	"healcationBackend/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	dsn := os.Getenv("DB")
	if dsn == "" {
		panic("Database DSN is not set")
	}
	fmt.Println("Connecting to database with DSN:", dsn)

	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	fmt.Println("Database connected")

	DB = connection
	fmt.Println("Starting migrations...")

	// Menjalankan auto migration
	err = connection.AutoMigrate(
		&models.History{},
		&models.SelectedAccomodation{},
		&models.SelectedPlace{},
		&models.User{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
	fmt.Println("Migrations completed")

	Seed()

	fmt.Println("Seed data inserted")
}
