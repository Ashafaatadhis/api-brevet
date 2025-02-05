package models

import (
	"time"
)

// Pertemuan model untuk table pertemuans
type Pertemuan struct {
	ID        int       `gorm:"primaryKey;type:int unsigned;"`                                         // Primary key
	Name      string    `gorm:"size:100;not null;uniqueIndex:idx_groupbatch_name,priority:2"`          // Nama pertemuan, unik per GroupBatch
	GrBatchID int       `gorm:"not null;type:int unsigned;uniqueIndex:idx_groupbatch_name,priority:1"` // Foreign key ke GroupBatch
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relasi belongs-to: Pertemuan memiliki satu GroupBatch
	GroupBatch *GroupBatch `gorm:"foreignKey:GrBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
