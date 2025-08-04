package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds all configuration for the application
type Config struct {
	ServerPort    string
	MongoURI      string
	MongoDatabase string
	JWTSecretKey  string
	JWTExpiration string
}

// LoadConfig reads configuration from environment variables.
// It will optionally load a .env file if it exists, but will not
// fail if it doesn't. This makes it work seamlessly in both
// local development and Docker environments.
func LoadConfig() (config Config, err error) {
	// Attempt to load .env file but don't treat it as a fatal error.
	// This is useful for local development.
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Could not find .env file. Using environment variables instead.")
	}

	// Read configuration from environment variables
	config = Config{
		ServerPort:    os.Getenv("SERVER_PORT"),
		MongoURI:      os.Getenv("MONGO_URI"),
		MongoDatabase: os.Getenv("MONGO_DATABASE"),
		JWTSecretKey:  os.Getenv("JWT_SECRET_KEY"),
		JWTExpiration: os.Getenv("JWT_EXPIRATION_HOURS"),
	}
	return config, nil // No error is returned from this function anymore
}
