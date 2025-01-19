package models

import "time"

// Profile adalah model untuk profiles
type Profile struct {
	ID         int       `gorm:"primaryKey;autoIncrement;type:int unsigned;"`
	GolonganID *int      `gorm:"omitempty;type:int unsigned;"`
	UserID     *int      `gorm:"omitempty;type:int unsigned;"` // Pointer untuk mendukung nilai NULL
	NIM        *string   `gorm:"size:15;omitempty"`
	BuktiNIM   *string   `gorm:"size:100;omitempty"`
	NIK        *string   `gorm:"size:20;omitempty"`
	BuktiNIK   *string   `gorm:"size:100;omitempty"`
	Institusi  string    `gorm:"size:100;not null"`
	Asal       string    `gorm:"size:100;not null"`
	TglLahir   time.Time `gorm:"not null"`
	Alamat     string    `gorm:"size:255;not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User     *User            `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Golongan KategoriGolongan `gorm:"foreignKey:GolonganID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
