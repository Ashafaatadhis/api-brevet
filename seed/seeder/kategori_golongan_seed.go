package seeder

import (
	"fmt"
	"new-brevet-be/config"
	"new-brevet-be/models"
)

// KategoriGolonganSeed untuk seed model Kategori_golongans
func KategoriGolonganSeed() {
	db := config.DB

	jenisKursus := []models.KategoriGolongan{
		{Name: "Umum"},
		{Name: "Mahasiswa"},
	}

	// Seed Data
	for _, jKurs := range jenisKursus {
		if err := db.FirstOrCreate(&jKurs, models.KategoriGolongan{Name: jKurs.Name}).Error; err != nil {
			fmt.Printf("Error seeding data: %v\n", err)
		}
	}

}
