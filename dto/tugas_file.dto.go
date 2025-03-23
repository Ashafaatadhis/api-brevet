package dto

import "time"

// TugasFileResponse dto response
type TugasFileResponse struct {
	ID      int    `json:"id"`
	TugasID int    `json:"tugas_id"` // Relasi ke Tugas
	FileURL string `json:"file_url"` // Lokasi file

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
