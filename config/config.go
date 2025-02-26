package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JWTSecret []byte
var Port string

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using default settings")
	}

	JWTSecret = []byte(getEnv("JWT_SECRET", "default"))
	Port = string(getEnv("PORT", "3000"))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
