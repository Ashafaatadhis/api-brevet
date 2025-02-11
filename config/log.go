package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// InitLogger mengatur konfigurasi global Logrus
func InitLogger() {
	// Set output ke stdout
	log.SetOutput(os.Stdout)

	// Paksa format JSON agar konsisten
	log.SetFormatter(&log.JSONFormatter{})

	// Set level log (opsional, bisa diubah sesuai kebutuhan)
	log.SetLevel(log.InfoLevel)
}
