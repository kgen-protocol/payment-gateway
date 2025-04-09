package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the environment configuration values used by the app.
type Config struct {
	Port         string
	MongoURI     string
	JwtSecret    string
	ClientID     string
	ClientSecret string
	GrantType    string
}

var instance *Config

// getEnvWithDefault returns the value of the environment variable identified by key,
// or returns defaultValue if the key is not set.
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Init loads the environment variables from .env (if present), then initializes the
// Config singleton using only the required environment variables.
// Note: You should only use the singleton from the service layer instead of calling os.Getenv directly.
func Init() {
	// Load .env file (for local development) if available.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize the config instance with required environment variables.
	instance = &Config{
		Port:         getEnvWithDefault("PORT", "8080"),
		MongoURI:     getEnvWithDefault("MONGO_URI", ""),
		JwtSecret:    getEnvWithDefault("JWT_SECRET", ""),
		ClientID:     getEnvWithDefault("CLIENT_ID", ""),
		ClientSecret: getEnvWithDefault("CLIENT_SECRET", ""),
		GrantType:    getEnvWithDefault("GRANT_TYPE", "client_credentials"),
	}
}

// GetConfig returns the singleton configuration instance.
// It lazily initializes the configuration if it has not been created yet.
func GetConfig() *Config {
	if instance == nil {
		Init()
	}
	return instance
}
