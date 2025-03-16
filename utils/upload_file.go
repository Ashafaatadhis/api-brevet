package utils

import (
	"fmt"
	"log"
	"mime/multipart"

	"time"

	"github.com/gofiber/fiber/v2"
)

// UploadFileHandler fungsi untuk menghandling upload file
func UploadFileHandler(c *fiber.Ctx, formFile *multipart.FileHeader, path *string) (*string, error) {

	// Cek ukuran file
	if formFile.Size > MaxFileSize {
		return nil, fiber.NewError(fiber.StatusBadRequest, "File size exceeds the 20MB limit")
	}

	// Cek MIME type gambar
	file, err := formFile.Open()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid file")
	}
	defer file.Close()

	// Mendapatkan tipe mime file
	mimeType, err := GetFileMimeType(file)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Unable to detect MIME type")
	}

	// Validasi apakah file adalah gambar
	if !IsValidDocument(mimeType) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid image format")
	}

	// Pastikan direktori upload tersedia
	if err := MakeDirectoryUploads(path); err != nil {
		// Lemparkan error jika direktori gagal dibuat

		log.Printf("failed to prepare upload directory: %w", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Failed to prepare upload directory")
	}

	uniqueFilename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), formFile.Filename)

	if path != nil {
		uniqueFilename = *path + "/" + uniqueFilename
	}

	uploadPath := "./public/uploads/" + uniqueFilename

	if err := c.SaveFile(formFile, uploadPath); err != nil {
		// Lemparkan error jika gagal menyimpan file
		log.Printf("failed to save image: %w", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Failed to save image")
	}

	return &uniqueFilename, nil // Tidak ada error
}

// UploadFile fungsi untuk upload file yang siap pakai dibanding uploadFileHandler
func UploadFile(c *fiber.Ctx, fieldName string, path string) (*string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		if err == fiber.ErrUnprocessableEntity {
			return nil, nil // Tidak ada file yang diupload, ini bukan error
		}
		return nil, err // Error saat membaca file
	}

	// Mendapatkan file path yang berupa string
	filePath, err := UploadFileHandler(c, file, &path)
	if err != nil {
		return nil, err
	}

	// Mengubah filePath (string) menjadi pointer string (*string)
	return filePath, nil
}
