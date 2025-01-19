package models

import "time"

// User adalah representasi tabel users di database
type User struct {
	ID        int       `gorm:"primaryKey;type:int unsigned;"` // ID user, primary key
	Name      string    `gorm:"size:100;not null"`             // Nama user
	Username  string    `gorm:"size:100;not null;uniqueIndex"` // Username user, harus unik
	Nohp      string    `gorm:"size:100;not null;uniqueIndex"` // No HP user, harus unik
	Avatar    string    `gorm:"size:255;not null"`             // Avatar URL user
	RoleID    int       `gorm:"type:int unsigned;not null"`    // ID role, tidak boleh null
	Email     string    `gorm:"uniqueIndex;size:100;not null"` // Email user, harus unik
	Password  string    `gorm:"size:255;not null"`             // Password user
	CreatedAt time.Time `gorm:"autoCreateTime"`                // Tanggal dibuat
	UpdatedAt time.Time `gorm:"autoUpdateTime"`                // Tanggal diperbarui

	Profile Profile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`               // Relasi dengan Profile
	Role    Role    `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Relasi dengan Role
}
