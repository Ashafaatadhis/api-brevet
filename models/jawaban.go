package models

import (
	"time"
)

// Jawaban models untuk table jawabans
type Jawaban struct {
	ID        int       `gorm:"primaryKey;type:int unsigned;"`
	TugasID   int       `gorm:"not null;type:int unsigned;uniqueIndex:idx_tugas_user"`
	UserID    int       `gorm:"not null;type:int unsigned;uniqueIndex:idx_tugas_user"` // One User, One Answer Per Tugas
	Answer    string    `gorm:"type:text;not null;"`
	Score     *int      `gorm:"type:int;default:null;"`
	Feedback  *string   `gorm:"type:text;default:null;"`
	IsLate    bool      `gorm:"not null;default:false;"`
	IsGraded  bool      `gorm:"not null;default:false;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	JawabanFile []JawabanFile `gorm:"foreignKey:JawabanID;"`

	Tugas Tugas `gorm:"foreignKey:TugasID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User  User  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
