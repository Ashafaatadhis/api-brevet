package models

import "time"

// TokenBlacklist adalah representasi tabel token_blacklists di database
type TokenBlacklist struct {
	ID        int       `gorm:"primaryKey;type:int unsigned;"`
	Token     string    `gorm:"not null"`
	ExpiredAt time.Time `gorm:"not null"`
}
