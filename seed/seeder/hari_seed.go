package seeder

import (
	"fmt"
	"new-brevet-be/config"
	"new-brevet-be/models"
)

// HariSeed untuk seed model Haris
func HariSeed() {
	db := config.DB

	hariList := []models.Hari{
		{Nama: "Senin"},
		{Nama: "Selasa"},
		{Nama: "Rabu"},
		{Nama: "Kamis"},
		{Nama: "Jumat"},
		{Nama: "Sabtu"},
		{Nama: "Minggu"},
	}

	// Seed Data
	for _, hari := range hariList {
		if err := db.FirstOrCreate(&hari, models.Hari{Nama: hari.Nama}).Error; err != nil {
			fmt.Printf("Error seeding data: %v\n", err)
		}
	}

}
