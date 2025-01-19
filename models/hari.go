package models

import "time"

// Hari adalah model untuk hari hari
type Hari struct {
	ID        int       `gorm:"primaryKey;type:int unsigned;"`
	Nama      string    `gorm:"size:50;not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Kursus []Kursus `gorm:"many2many:group_hr_kursus;"`
}
