package tasks

import (
	"fmt"
	"log"
	"new-brevet-be/config"
	"new-brevet-be/models"
	"os"
	"strconv"
	"time"
)

// CleanupExpiredTokens adalah fungsi untuk membersihkan token yang sudah expired di table blacklisttoken
func CleanupExpiredTokens() {
	var timeClean = os.Getenv("GOROUTINE_CLEAN_TOKEN")
	db := config.DB
	expiryInHours, err := strconv.Atoi(timeClean)
	if err != nil {
		fmt.Println("Error parsing TOKEN_EXPIRY:", err)
		return
	}

	ticker := time.NewTicker(time.Duration(expiryInHours) * time.Hour) // Jalankan setiap 1 jam
	defer ticker.Stop()

	for {
		<-ticker.C

		// Hapus token yang expired
		if err := db.Where("expired_at < ?", time.Now()).Delete(&models.TokenBlacklist{}).Error; err != nil {
			log.Println("Failed to clean up expired tokens:", err)
		} else {
			log.Println("Expired tokens cleaned up successfully")
		}
	}
}
