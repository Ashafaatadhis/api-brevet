package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv meload env variables
func LoadEnv() {
	// Cek path kerja saat unit test
	dir, _ := os.Getwd()
	log.Println("Current Working Directory:", dir)

	// Load dari beberapa kemungkinan path
	envPaths := []string{
		".env",       // Jika berjalan dari root project
		"../.env",    // Jika berjalan dari folder tests
		"../../.env", // Jika test dijalankan lebih dalam
	}

	var err error
	for _, path := range envPaths {
		err = godotenv.Load(path)
		if err == nil {
			log.Println("✅ Loaded env file:", path)
			return
		}
	}

	log.Fatal("❌ Error loading .env file")
}
