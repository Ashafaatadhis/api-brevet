package models

import "time"

// Materi models untuk table materis
type Materi struct {
	ID          int       `gorm:"primaryKey;type:int unsigned;"`
	PertemuanID int       `gorm:"not null;type:int unsigned;"` // Relasi ke pertemuan
	Title       string    `gorm:"size:255"`
	Content     string    `gorm:"type:text"` // Bisa berupa teks
	FileURL     string    `gorm:"size:255"`  // Jika ada file (PDF, video, dsb)
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Pertemuan Pertemuan `gorm:"foreignKey:PertemuanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
