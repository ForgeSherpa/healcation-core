package config

import (
	"log"
	"os"
	"slices"

	"github.com/joho/godotenv"
)

var (
	GeminiAPIKey     string
	GoogleAPIKey     string
	GoogleAPI_CX     string
	TursoDatabaseURL string
	TursoAuthToken   string
	AppEnv           string
	IsStaging        bool
	IsProduction     bool
	// IsGeminiEnabled is a flag to enable/disable Gemini API usage (enum: "1" or "0")
	IsGeminiEnabled bool
)

func loadAppEnv() {
	AppEnv = os.Getenv("APP_ENV")

	// do not load .env file in production
	if AppEnv == "production" {
		return
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func init() {
	loadAppEnv()

	GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
	GoogleAPI_CX = os.Getenv("GOOGLE_API_CX")
	GoogleAPIKey = os.Getenv("GOOGLE_API_KEY")

	TursoDatabaseURL = os.Getenv("TURSO_DATABASE_URL")
	TursoAuthToken = os.Getenv("TURSO_AUTH_TOKEN")

	AppEnv = os.Getenv("APP_ENV")

	if !slices.Contains([]string{"staging", "production"}, AppEnv) {
		log.Fatal("environment variable APP_ENV must be 'staging' or 'production'")
	}

	IsStaging = AppEnv == "staging"
	IsProduction = AppEnv == "production"

	// no need to cast, as this check already returs a boolean
	IsGeminiEnabled = os.Getenv("IS_GEMINI_ENABLED") == "1"
}
