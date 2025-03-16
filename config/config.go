package config

import (
	"fmt"
	"log"
	"os"

	// untuk menjalankan driver
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

// DB variabel yang merepresentasikan db dari gorm DB
var DB *gorm.DB

// InitDB untuk koneksi ke db
func InitDB() {
	var err error

	// Mengambil DSN dari environment variable
	dsn := os.Getenv("DATABASE_URL")

	// Menggunakan gorm.Open dengan mysql.Open(dsn) untuk GORM v2
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database: ", err)
	}

	fmt.Println("Database connected")
}

// TestDB variabel yang merepresentasikan db dari gorm DB
var TestDB *gorm.DB

// InitTestDB untuk koneksi ke db testing
func InitTestDB() {

	var err error

	// Mengambil DSN dari environment variable
	dsn := os.Getenv("DATABASE_URL_TEST")

	// Menggunakan gorm.Open dengan mysql.Open(dsn) untuk GORM v2
	TestDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database: ", err)
	}

	fmt.Println("Database connected")
}
