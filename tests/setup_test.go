package tests

import (
	"log"
	"os"
	"testing"

	"new-brevet-be/config"
	"new-brevet-be/handlers"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var testApp *fiber.App

// 🔹 Setup Database Sekali untuk Semua Test
func setupTestDB() *gorm.DB {
	config.InitTestDB()
	db := config.TestDB

	// 🔥 Drop & Migrate ulang semua tabel
	db.Migrator().DropTable(&models.User{}, &models.Profile{}, &models.Role{}, &models.KategoriGolongan{}, &models.TokenBlacklist{},
		&models.Batch{}, &models.JenisKursus{}, &models.KelasKursus{}, &models.GroupBatch{}, &models.Hari{}, &models.Kursus{}, &models.Category{})
	db.AutoMigrate(&models.User{}, &models.Profile{}, &models.Role{}, &models.KategoriGolongan{}, &models.TokenBlacklist{},
		&models.Batch{}, &models.JenisKursus{}, &models.KelasKursus{}, &models.GroupBatch{}, &models.Hari{}, &models.Kursus{}, &models.Category{})

	// ✅ Seed Data (Masukkan Data Penting)
	seedDatabase(db)

	log.Println("✅ Database mock berhasil diinisialisasi!")
	return db
}

// 🔹 Seed Data Penting (Tanpa Query Raw)
func seedDatabase(db *gorm.DB) {
	roles := []models.Role{
		{Name: "admin"},
		{Name: "helpdesk"},
		{Name: "guru"},
		{Name: "siswa"},
	}
	kelasKursus := []models.KelasKursus{
		{Name: "Reguler"},
		{Name: "Eksekutif"},
	}
	jenisKursus := []models.JenisKursus{
		{Name: "Online"},
		{Name: "Offline"},
	}
	categories := []models.Category{
		{Name: "Kursus"},
		{Name: "Seminar/Workshop"},
		{Name: "Sertifikasi"},
	}
	haris := []models.Hari{
		{Nama: "Senin"},
		{Nama: "Selasa"},
		{Nama: "Rabu"},
		{Nama: "Kamis"},
		{Nama: "Jumat"},
		{Nama: "Sabtu"},
		{Nama: "Minggu"},
	}

	db.Create(&roles)
	db.Create(&kelasKursus)
	db.Create(&jenisKursus)
	db.Create(&haris)
	db.Create(&categories)

	log.Println("✅ Data role, kelas, dan jenis kursus berhasil di-seed!")
}

// 🔹 Bersihkan Database Setelah Tiap Test
func cleanupDatabase(db *gorm.DB) {
	log.Print("TEST CLEAN")

	// Nonaktifkan sementara FOREIGN_KEY_CHECKS untuk menghindari constraint error
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// Truncate semua tabel
	db.Exec("TRUNCATE TABLE users")
	db.Exec("TRUNCATE TABLE profiles")

	// db.Exec("TRUNCATE TABLE kategori_golongan")
	db.Exec("TRUNCATE TABLE token_blacklists")
	db.Exec("TRUNCATE TABLE batches")
	db.Exec("TRUNCATE TABLE kursus")

	db.Exec("TRUNCATE TABLE group_batches")

	// Aktifkan kembali FOREIGN_KEY_CHECKS setelah truncate
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	log.Println("🧹 Database dibersihkan setelah test selesai.")
}
func cleanupDatabaseAll(db *gorm.DB) {
	log.Print("TEST CLEAN")

	// Nonaktifkan sementara FOREIGN_KEY_CHECKS untuk menghindari constraint error
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// Truncate semua tabel
	db.Exec("TRUNCATE TABLE users")
	db.Exec("TRUNCATE TABLE profiles")
	db.Exec("TRUNCATE TABLE roles")
	// db.Exec("TRUNCATE TABLE kategori_golongan")
	db.Exec("TRUNCATE TABLE token_blacklists")
	db.Exec("TRUNCATE TABLE batches")
	db.Exec("TRUNCATE TABLE jenis_kursus")
	db.Exec("TRUNCATE TABLE kelas_kursus")
	db.Exec("TRUNCATE TABLE group_batches")
	db.Exec("TRUNCATE TABLE kursus")
	db.Exec("TRUNCATE TABLE categories")

	// Aktifkan kembali FOREIGN_KEY_CHECKS setelah truncate
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	log.Println("🧹 Database dibersihkan setelah test selesai.")
}

// 🔹 Setup Fiber App (Gunakan Database yang Sudah Dibuat)
func setupApp() *fiber.App {
	app := fiber.New()

	app.Get("/me", middlewares.AuthMiddleware(), handlers.Me)
	app.Post("/register", validation.Validate[validation.UserRegister](), middlewares.UserUniqueCheck[validation.UserRegister], handlers.Register())
	app.Post("/login", validation.Validate[validation.UserLogin](), handlers.Login())
	app.Delete("/logout", middlewares.AuthMiddleware(), handlers.Logout())

	app.Get("/batch", handlers.GetBatch)
	app.Get("/batch/:id", handlers.GetDetailBatch)
	app.Post("/batch", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}), validation.Validate[validation.PostBatch](), middlewares.BatchUniqueCheck[validation.PostBatch], handlers.PostBatch)
	app.Put("/batch/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}), validation.Validate[validation.PostBatch](), middlewares.BatchUniqueCheck[validation.PostBatch], handlers.UpdateBatch)
	app.Delete("/batch/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}), handlers.DeleteBatch)

	app.Get("/kursus", handlers.GetKursus)
	app.Get("/kursus/:id", handlers.GetDetailKursus)
	app.Post("/kursus", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.PostKursus](), middlewares.KursusUniqueCheck[validation.PostKursus],
		handlers.PostKursus)
	app.Put("/kursus/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.PostKursus](), middlewares.KursusUniqueCheck[validation.PostKursus],
		handlers.UpdateKursus)
	app.Delete("/kursus/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		handlers.DeleteKursus)

	// batch-mapping
	app.Get("/batch-mapping", handlers.GetAllBatchMappping)
	app.Get("/batch-mapping/:id", handlers.GetDetailBatchMappping)
	app.Post("/batch-mapping", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.CreateBatchMapping](),
		handlers.CreateBatchMapping)
	app.Put("/batch-mapping/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.CreateBatchMapping](),
		handlers.EditBatchMapping)
	app.Delete("/batch-mapping/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		handlers.DeleteBatchMapping)

	return app
}

// 🔹 TestMain: Jalankan Sekali Sebelum Semua Test
func TestMain(m *testing.M) {
	log.Println("🔹 Setup test environment...")
	config.LoadEnv()

	testDB = setupTestDB()
	config.DB = testDB

	testApp = setupApp()

	// Jalankan semua test
	exitCode := m.Run()

	// Bersihkan database setelah semua test selesai
	cleanupDatabaseAll(testDB)

	log.Println("✅ Semua test selesai, database dibersihkan.")
	os.Exit(exitCode)
}
