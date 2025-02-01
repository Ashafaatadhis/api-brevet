package main

import (
	"log"

	"new-brevet-be/config"
	"new-brevet-be/models"
)

func main() {
	config.LoadEnv()
	// Hubungkan ke database
	config.InitDB()

	// Panggil fungsi untuk migrasi tabel
	migrate()

	log.Println("Migration completed successfully")
}

// Fungsi untuk melakukan migrasi
func migrate() {
	db := config.DB
	// Migrasi untuk model Role
	// db.DropTable(&models.User{}, &models.Role{}, &models.TokenBlacklist{})
	// Menambahkan foreign key untuk kolom RoleID di tabel User
	// db.AutoMigrate(&models.Batch{}, &models.TokenBlacklist{}, &models.Category{}, &models.GroupBatch{}, &models.Hari{},
	// 	&models.JenisKursus{}, &models.KategoriGolongan{}, &models.KelasKursus{}, &models.Kursus{}, &models.Profile{},
	// 	&models.Purchase{}, &models.Role{},
	// 	&models.StatusPayment{}, &models.User{})
	db.Migrator().DropTable(&models.Kursus{})
	db.AutoMigrate(&models.Kursus{}, &models.User{}, &models.JenisKursus{}, &models.KelasKursus{}, &models.Category{}, &models.GroupBatch{})
}
