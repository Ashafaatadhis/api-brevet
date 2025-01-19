package seeder

import (
	"fmt"
	"new-brevet-be/config"
	"new-brevet-be/models"
)

// JenisKursusSeed untuk seed model Jenis Kursus
func JenisKursusSeed() {
	db := config.DB

	jenisKursus := []models.KelasKursus{
		{Name: "Reguler"},
		{Name: "Eksekutif"},
	}

	// Seed Data
	for _, jKurs := range jenisKursus {
		if err := db.FirstOrCreate(&jKurs, models.KelasKursus{Name: jKurs.Name}).Error; err != nil {
			fmt.Printf("Error seeding data: %v\n", err)
		}
	}

}
