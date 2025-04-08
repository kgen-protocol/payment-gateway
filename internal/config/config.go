package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file if present
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}
