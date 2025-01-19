package seeder

import (
	"log"

	"new-brevet-be/config"
	"new-brevet-be/models"
)

// RolesSeed untuk tabel Role
func RolesSeed() {
	db := config.DB
	adminRole := models.Role{Name: "admin"}
	helpDeskRole := models.Role{Name: "helpdesk"}
	guruRole := models.Role{Name: "guru"}
	siswaRole := models.Role{Name: "siswa"}

	// Memeriksa apakah role sudah ada, jika tidak, buat
	db.FirstOrCreate(&adminRole, &adminRole)

	db.FirstOrCreate(&helpDeskRole, &helpDeskRole)
	db.FirstOrCreate(&guruRole, &guruRole)
	db.FirstOrCreate(&siswaRole, &siswaRole)
	log.Println("Roles seeded successfully")
}
