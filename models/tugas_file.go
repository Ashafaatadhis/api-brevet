package models

import "time"

// TugasFile menyimpan file yang diupload untuk setiap tugas
type TugasFile struct {
	ID      int    `gorm:"primaryKey;type:int unsigned;"`
	TugasID *int   `gorm:"default:null;type:int unsigned;"`
	FileURL string `gorm:"type:varchar(255);not null;"` // Lokasi file

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Tugas     Tugas     `gorm:"foreignKey:TugasID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
