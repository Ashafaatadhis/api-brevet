package models

import "time"

// Materi models untuk table materis
type Materi struct {
	ID          int       `gorm:"primaryKey;type:int unsigned;"`
	PertemuanID int       `json:"pertemuan_id"` // Relasi ke pertemuan
	Title       string    `json:"title"`
	Content     string    `json:"content"`  // Bisa berupa teks
	FileURL     string    `json:"file_url"` // Jika ada file (PDF, video, dsb)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Pertemuan Pertemuan `gorm:"foreignKey:PertemuanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
