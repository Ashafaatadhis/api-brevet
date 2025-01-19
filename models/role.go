package models

// Role adalah representasi tabel roles di database
type Role struct {
	ID int `gorm:"primaryKey;type:int unsigned"` // Pastikan ini sesuai
	// ID role, primary key
	Name string `gorm:"size:100;not null;uniqueIndex"` // Nama role, harus unik

}
