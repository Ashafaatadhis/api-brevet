package models

import "time"

// Jawaban models untuk table jawabans
type Jawaban struct {
	ID          int       `gorm:"primaryKey;type:int unsigned;"`
	PertemuanID int       `json:"pertemuan_id"`
	UserID      int       `json:"user_id"`
	Answer      string    `json:"answer"`   // Untuk jawaban teks
	FileURL     string    `json:"file_url"` // Untuk jawaban file
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Pertemuan Pertemuan `gorm:"foreignKey:PertemuanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
