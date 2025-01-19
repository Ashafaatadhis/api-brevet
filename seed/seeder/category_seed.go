package seeder

import (
	"fmt"
	"new-brevet-be/config"
	"new-brevet-be/models"
)

// CategorySeed untuk seed model categories
func CategorySeed() {
	db := config.DB

	jenisKursus := []models.Category{
		{Name: "Kursus"},
		{Name: "Seminar/Workshop"},
		{Name: "Sertifikasi"},
	}

	// Seed Data
	for _, jKurs := range jenisKursus {
		if err := db.FirstOrCreate(&jKurs, models.Category{Name: jKurs.Name}).Error; err != nil {
			fmt.Printf("Error seeding data: %v\n", err)
		}
	}

}
