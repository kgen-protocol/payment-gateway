package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	// Pinelabs specific credentials
	PinelabsClientID     string
	PinelabsClientSecret string
	PinelabsGrantType    string
	PinelabsTokenURL     string
	PinelabsOrderURL     string
	PinelabsGetOrderURL  string

	// Redis configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	RedisAddr     string
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

func parseEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return defaultVal
}

// Init loads the environment variables from .env (if present), then initializes the
// Config singleton using only the required environment variables.
// Note: You should only use the singleton from the service layer instead of calling os.Getenv directly.
func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	redisDB := 0
	if dbEnv := os.Getenv("REDIS_DB"); dbEnv != "" {
		fmt.Sscanf(dbEnv, "%d", &redisDB)
	}

	host := getEnvWithDefault("REDIS_HOST", "localhost")
	port := getEnvWithDefault("REDIS_PORT", "6379")

	instance = &Config{
		Port:         getEnvWithDefault("PORT", "8080"),
		MongoURI:     getEnvWithDefault("MONGO_URI", ""),
		JwtSecret:    getEnvWithDefault("JWT_SECRET", ""),
		ClientID:     getEnvWithDefault("CLIENT_ID", ""),
		ClientSecret: getEnvWithDefault("CLIENT_SECRET", ""),
		GrantType:    getEnvWithDefault("GRANT_TYPE", "client_credentials"),

		PinelabsClientID:     getEnvWithDefault("PINELABS_CLIENT_ID", ""),
		PinelabsClientSecret: getEnvWithDefault("PINELABS_CLIENT_SECRET", ""),
		PinelabsGrantType:    getEnvWithDefault("PINELABS_GRANT_TYPE", "client_credentials"),
		PinelabsTokenURL:     getEnvWithDefault("PINELABS_TOKEN_URL", ""),
		PinelabsOrderURL:     getEnvWithDefault("PINELABS_ORDER_URL", ""),
		PinelabsGetOrderURL:  getEnvWithDefault("PINELABS_GET_ORDER_URL", ""),

		RedisHost:     host,
		RedisPort:     port,
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		RedisAddr:     host + ":" + port,
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
