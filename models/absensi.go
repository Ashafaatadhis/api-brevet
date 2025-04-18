package models

import (
	"time"
)

// Absensi model
type Absensi struct {
	ID           int       `gorm:"primaryKey;type:int unsigned;"`
	UserID       int       `gorm:"not null;type:int unsigned;index:idx_user_batch_date,unique"`
	GroupBatchID int       `gorm:"not null;type:int unsigned;index:idx_user_batch_date,unique"`
	Tanggal      time.Time `gorm:"type:date;not null"`
	Status       string    `gorm:"type:enum('hadir','izin','sakit','alfa');default:'hadir'"`

	GuruID *int `gorm:"type:int unsigned"` // harus nullable
	Guru   User `gorm:"foreignKey:GuruID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User       User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GroupBatch GroupBatch `gorm:"foreignKey:GroupBatchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
