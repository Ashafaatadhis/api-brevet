package seeder

import (
	"fmt"
	"new-brevet-be/config"
	"new-brevet-be/models"
)

// StatusPaymentSeed untuk seed model status_payments
func StatusPaymentSeed() {
	db := config.DB

	jenisPayment := []models.StatusPayment{
		{Name: "Pending"},
		{Name: "Lunas"},
		{Name: "Gagal"},
	}

	// Seed Data
	for _, jPay := range jenisPayment {
		if err := db.FirstOrCreate(&jPay, models.StatusPayment{Name: jPay.Name}).Error; err != nil {
			fmt.Printf("Error seeding data: %v\n", err)
		}
	}

}
