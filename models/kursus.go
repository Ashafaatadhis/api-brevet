package models

import (
	"time"
)

// Kursus adalah representasi tabel kursus di database
type Kursus struct {
	ID               int       `gorm:"primaryKey;type:int unsigned;"`
	TeacherID        *int      `gorm:"size:255;omitempty;type:int unsigned;"`
	Judul            string    `gorm:"size:255;not null"`
	JenisID          int       `gorm:"not null;type:int unsigned;"`
	KelasID          int       `gorm:"not null;type:int unsigned;"`
	DeskripsiSingkat string    `gorm:"size:255"`
	Deskripsi        string    `gorm:"type:text"`
	Pembelajaran     string    `gorm:"type:text"`
	Diperoleh        string    `gorm:"type:text"`
	CategoryID       int       `gorm:"not null;type:int unsigned;"`
	ThumbnailKursus  string    `gorm:"size:255"`
	ThumbnailURL     string    `gorm:"size:255"`
	HargaAsli        float64   `gorm:"not null"`
	HargaDiskon      float64   `gorm:"not null"`
	StartDate        time.Time `gorm:"not null"`
	EndDate          time.Time `gorm:"not null"`
	StartTime        time.Time `gorm:"not null"`
	EndTime          time.Time `gorm:"not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`

	Teacher  *User       `gorm:"foreignKey:TeacherID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Jenis    JenisKursus `gorm:"foreignKey:JenisID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Kelas    KelasKursus `gorm:"foreignKey:KelasID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Category Category    `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	GroupBatches []*GroupBatch `gorm:"foreignKey:KursusID"`
	Hari         []Hari        `gorm:"many2many:group_hr_kursus;"`
}
