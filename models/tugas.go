package models

import "time"

// Tugas models untuk table tugas
type Tugas struct {
	ID          int       `gorm:"primaryKey;type:int unsigned;"`
	PertemuanID int       `gorm:"not null;type:int unsigned;"` // Relasi ke pertemuan
	Title       string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text;not null"`
	Type        string    `gorm:"type:ENUM('essay', 'file');not null"` // "essay" atau "file"
	StartAt     time.Time `gorm:"not null"`                            // Tugas mulai bisa dikerjakan
	EndAt       time.Time `gorm:"not null"`                            // Batas akhir pengumpulan tugas
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Jawaban []Jawaban `gorm:"foreignKey:TugasID;"`

	Pertemuan Pertemuan   `gorm:"foreignKey:PertemuanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TugasFile []TugasFile `gorm:"foreignKey:TugasID;"`
}
