package models

import "time"

// Tugas models untuk table tugas
type Tugas struct {
	ID          int       `gorm:"primaryKey;type:int unsigned;"`
	PertemuanID int       `json:"pertemuan_id"` // Relasi ke pertemuan
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // "essay" atau "file"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Pertemuan Pertemuan `gorm:"foreignKey:PertemuanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
