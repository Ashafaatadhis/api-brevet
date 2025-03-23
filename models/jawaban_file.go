package models

import "time"

// JawabanFile models untuk table jawaban_files
type JawabanFile struct {
	ID        int       `gorm:"primaryKey;type:int unsigned;"`
	JawabanID int       `gorm:"not null;type:int unsigned;"` // Wajib ada jawaban
	FileURL   string    `gorm:"type:text;not null;"`         // Wajib ada file URL
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Jawaban Jawaban `gorm:"foreignKey:JawabanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
