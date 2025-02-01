package seeder

import (
	"log"
	"math/rand"
	"new-brevet-be/config"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"time"

	"github.com/bxcodec/faker/v3"
)

// UserSeed untuk tabel User
func UserSeed() {
	db := config.DB

	// Seed 30 users
	for i := 0; i < 30; i++ {
		// Hash password
		hashedPassword, err := utils.HashPassword("password123")
		if err != nil {
			log.Println("Failed to hash password:", err)
			return
		}

		// Generate random role from 1 to 4
		roleID := rand.Intn(4) + 1

		// Create user with random data
		user := models.User{
			Username: faker.Username(),
			Name:     faker.Name(),
			Nohp:     faker.Phonenumber(),
			Email:    faker.Email(),
			RoleID:   roleID,
			Password: hashedPassword,
		}

		// Create user in the database
		if err := db.Create(&user).Scan(&user).Error; err != nil {
			log.Printf("Failed to create user: %v\n", err)
			return
		}

		// Create profile for the user
		profile := models.Profile{
			UserID:    &user.ID,
			Institusi: "Gundar",
			Asal:      "Indramayu",
			TglLahir:  generateRandomDate(),
			Alamat:    "Juntinyuat",
		}

		// Create profile in the database
		if err := db.Create(&profile).Error; err != nil {
			log.Printf("Failed to create profile for user %v: %v\n", user.Username, err)
			return
		}
	}

	log.Println("Users and profiles seeded successfully")
}

// generateRandomDate generates a random date of birth between 1970 and 2000
func generateRandomDate() time.Time {
	year := rand.Intn(30) + 1970           // Random year between 1970 and 2000
	month := time.Month(rand.Intn(12) + 1) // Random month between 1 and 12
	day := rand.Intn(28) + 1               // Random day between 1 and 28 (to avoid invalid dates)
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
