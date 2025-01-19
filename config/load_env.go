package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv meload env variables
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
