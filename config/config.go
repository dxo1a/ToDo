package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JWTSecret []byte

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using default settings")
	}

	JWTSecret = []byte(getEnv("JWT_SECRET", "default"))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
