package models

import "time"

// Batch adalah model untuk batch
type Batch struct {
	ID         int       `gorm:"primaryKey;type:int unsigned;"`
	Judul      string    `gorm:"size:100;not null;unique"`
	BukaBatch  time.Time `gorm:"not null"`
	TutupBatch time.Time `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	GroupBatches []*GroupBatch `gorm:"foreignKey:BatchID"`
}
