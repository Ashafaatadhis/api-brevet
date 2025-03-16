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

	db.AutoMigrate(&models.Pertemuan{}, &models.Materi{}, &models.Tugas{}, &models.Jawaban{}, &models.Blog{})

	// db.AutoMigrate(&models.GroupBatch{}, &models.Pertemuan{})
	// db.Migrator().DropTable(&models.Purchase{})
	// db.AutoMigrate(&models.Purchase{})

	// Migrasi untuk model Role
	// db.DropTable(&models.User{}, &models.Role{}, &models.TokenBlacklist{})
	// Menambahkan foreign key untuk kolom RoleID di tabel User
	// db.AutoMigrate(&models.Batch{}, &models.TokenBlacklist{}, &models.Category{}, &models.GroupBatch{}, &models.Hari{},
	// 	&models.JenisKursus{}, &models.KategoriGolongan{}, &models.KelasKursus{}, &models.Kursus{}, &models.Profile{},
	// 	&models.Purchase{}, &models.Role{},
	// 	&models.StatusPayment{}, &models.User{})
	// db.Migrator().DropColumn(&models.Kursus{}, "StartTime")
	// db.Migrator().DropColumn(&models.Kursus{}, "EndTime")

	// err := db.Migrator().AlterColumn(&models.Kursus{}, "StartTime")
	// if err != nil {
	// 	fmt.Println("Gagal mengubah StartTime:", err)
	// } else {
	// 	fmt.Println("StartTime berhasil diubah ke VARCHAR(8)")
	// }

	// err = db.Migrator().AlterColumn(&models.Kursus{}, "EndTime")
	// if err != nil {
	// 	fmt.Println("Gagal mengubah EndTime:", err)
	// } else {
	// 	fmt.Println("EndTime berhasil diubah ke VARCHAR(8)")
	// }
	// db.AutoMigrate(&models.Kursus{}) // Buat ulang kolom dengan tipe baru
	// db.Migrator().DropTable(&models.Kursus{})
	// db.AutoMigrate(&models.Kursus{}, &models.User{}, &models.JenisKursus{}, &models.KelasKursus{}, &models.Category{}, &models.GroupBatch{})
}
