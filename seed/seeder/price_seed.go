package seeder

import (
	"fmt"
	"new-brevet-be/config"
	"new-brevet-be/models"
)

// PriceSeed untuk seed model Price
func PriceSeed() {
	db := config.DB

	jenisKursus := []models.Price{
		{GolonganID: 1, Harga: 1000000},
		{GolonganID: 2, Harga: 750000},
	}

	// Seed Data
	for _, jKurs := range jenisKursus {
		if err := db.FirstOrCreate(&jKurs, models.Price{GolonganID: jKurs.GolonganID, Harga: jKurs.Harga}).Error; err != nil {
			fmt.Printf("Error seeding data: %v\n", err)
		}
	}

}
