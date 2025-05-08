package config

import (
	"log"
	"os"
	"strconv"

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
	IsGeminiEnabled  bool
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
	GoogleAPI_CX = os.Getenv("GOOGLE_API_CX")
	GoogleAPIKey = os.Getenv("GOOGLE_API_KEY")

	TursoDatabaseURL = os.Getenv("TURSO_DATABASE_URL")
	TursoAuthToken = os.Getenv("TURSO_AUTH_TOKEN")

	AppEnv = os.Getenv("APP_ENV")
	if AppEnv != "staging" && AppEnv != "production" {
		log.Fatal("environment variable APP_ENV must be 'staging' or 'production'")
	}
	IsStaging = (AppEnv == "staging")
	IsProduction = (AppEnv == "production")

	if v, err := strconv.ParseBool(os.Getenv("IS_GEMINI_ENABLED")); err == nil {
		IsGeminiEnabled = v
	} else {
		IsGeminiEnabled = false
	}
}
