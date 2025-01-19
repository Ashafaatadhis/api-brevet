package utils

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// GenerateUniqueFilename buat generate unik filename
func GenerateUniqueFilename(originalFilename string) string {
	// Ambil ekstensi file asli
	ext := filepath.Ext(originalFilename)

	// Buat UUID dan timestamp
	uuidPart := uuid.New().String()
	timestampPart := time.Now().Format("20060102150405")

	// Gabungkan menjadi nama file unik
	return fmt.Sprintf("%s_%s%s", uuidPart, timestampPart, ext)
}
