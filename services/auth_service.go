package services

import (
	"new-brevet-be/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CheckTokenBlacklist adalah fungsi untuk mencari token apakah ada di blacklist atau tidak
func CheckTokenBlacklist(db *gorm.DB, tokenString string, blacklist *models.TokenBlacklist) error {
	// Cek apakah token ada di blacklist
	err := db.Where("token = ?", tokenString).First(&blacklist).Error

	// Jika token ditemukan di blacklist, lempar error Unauthorized
	if err == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
	}

	// Jika ada error selain ErrRecordNotFound, lempar error Unauthorized
	if err != gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
	}

	// Jika token tidak ditemukan di blacklist, return nil untuk melanjutkan proses
	return nil
}
