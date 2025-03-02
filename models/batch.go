package models

import "time"

// Batch adalah model untuk batch
type Batch struct {
	ID         int       `gorm:"primaryKey;type:int unsigned;"`
	Judul      string    `gorm:"size:100;not null;unique"`
	JenisID    int       `gorm:"not null;type:int unsigned;"`
	KelasID    int       `gorm:"not null;type:int unsigned;"`
	BukaBatch  time.Time `gorm:"not null"`
	TutupBatch time.Time `gorm:"not null"`
	Kuota      int       `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	Jenis JenisKursus `gorm:"foreignKey:JenisID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Kelas KelasKursus `gorm:"foreignKey:KelasID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	GroupBatches []*GroupBatch `gorm:"foreignKey:BatchID"`
}
