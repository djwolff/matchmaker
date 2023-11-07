package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Setup() error {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
		return err
	}

	if os.Getenv("ENV") == "dev" {
		os.Setenv("APP_URL", "http://localhost:8080")
	}
	return nil
}
