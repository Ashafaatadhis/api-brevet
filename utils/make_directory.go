package utils

import (
	"fmt"
	"log"
	"os"
)

// MakeDirectoryUploads untuk membuat direktory di dalam folder public
func MakeDirectoryUploads(path *string) error {
	uploadDir := "./public/uploads/"
	if path != nil {
		uploadDir = uploadDir + *path
	}
	// Periksa apakah direktori sudah ada
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		// Buat direktori jika belum ada
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			log.Printf("Failed to create directory: %v", err)
			return fmt.Errorf("failed to prepare upload directory: %w", err) // Melempar error
		}
	}

	// Jika sukses, kembalikan nil
	return nil
}
