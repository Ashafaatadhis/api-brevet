package utils

import (
	"io"
	"net/http"
)

// MaxFileSize ukuran file 20MB
const MaxFileSize = 20 * 1024 * 1024 // 20MB

// IsValidImage untuk validasi tipe gambar
func IsValidImage(mimeType string) bool {
	// Validasi tipe mime file yang diperbolehkan
	allowedMimeTypes := []string{
		"image/jpeg", "image/png", "image/gif", "image/webp", "image/jpg",
	}

	for _, validType := range allowedMimeTypes {
		if mimeType == validType {
			return true
		}
	}
	return false
}

// GetFileMimeType untuk mendapatkan MIME type dari file
func GetFileMimeType(file io.Reader) (string, error) {
	buffer := make([]byte, 512) // Membaca sebagian kecil file untuk deteksi MIME
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buffer), nil
}
