package seeder

import (
	"log"

	"new-brevet-be/config"
	"new-brevet-be/models"
	"new-brevet-be/utils"
)

// UserSeed untuk tabel User
func UserSeed() {
	db := config.DB
	hashedPassword, err := utils.HashPassword("helpdesk123")
	if err != nil {
		log.Println("Failed to hash password:", err)
	}
	admin := models.User{Username: "helpdesk", Name: "helpdesk", Nohp: "082131", RoleID: 2, Email: "helpdesk@gmail.com", Password: hashedPassword}

	// Memeriksa apakah role sudah ada, jika tidak, buat

	db.Create(&admin)
	log.Println("User seeded successfully")
}
