package utils

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
)

// ParseMultipartForm menangani parsing multipart form dan mengambil file berdasarkan fieldName.
func ParseMultipartForm(c *fiber.Ctx, fieldName string) ([]*multipart.FileHeader, error) {
	// Ambil form data dari request
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	// Ambil file berdasarkan fieldName
	files, ok := form.File[fieldName]
	if !ok {
		return nil, nil // Kembalikan nil jika tidak ada file
	}

	return files, nil
}
