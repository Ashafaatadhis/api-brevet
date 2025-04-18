package dto

import (
	"time"
)

// AbsensiResponse struct untuk response khusus menangani data batch
type AbsensiResponse struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	GroupBatchID int       `json:"group_batch_id"`
	Tanggal      time.Time `json:"tanggal"`
	Status       string    `json:"status"`

	GuruID *int `json:"guru_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AbsensiRequest struct untuk request khusus menangani data absensi
type AbsensiRequest struct {
	UserID       string    `validate:"required,min=3,max=100" json:"user_id"`
	GroupBatchID time.Time `validate:"required" json:"group_batch_id"`
	Tanggal      time.Time `validate:"required" json:"tanggal"` // Tanggal tanpa jam
	Status       string    `validate:"required" json:"status"`
	GuruID       *int      `validate:"required" json:"guru_id"`
}
