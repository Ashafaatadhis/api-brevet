package models

import "time"

// Purchase adalah model untuk purchases
type Purchase struct {
	ID              int       `gorm:"primaryKey;autoIncrement;type:int unsigned;"`
	GrBatchID       int       `gorm:"not null;type:int unsigned;"`
	StatusPaymentID int       `gorm:"not null;type:int unsigned;"`
	JenisKursusID   int       `gorm:"not null;type:int unsigned;"`
	UserID          *int      `gorm:"omitempty;type:int unsigned;"`
	URLConfirm      *string   `gorm:"size:255;omitempty"`
	BuktiBayar      *string   `gorm:"size:100;omitempty"`
	PriceID         int       `gorm:"not null;type:int unsigned;"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`

	Price         Price         `gorm:"foreignKey:PriceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	JenisKursus   JenisKursus   `gorm:"foreignKey:JenisKursusID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GroupBatches  *GroupBatch   `gorm:"foreignKey:GrBatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User          *User         `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StatusPayment StatusPayment `gorm:"foreignKey:StatusPaymentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
