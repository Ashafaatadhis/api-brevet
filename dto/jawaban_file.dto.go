package dto

import "time"

// JawabanFileResponse dto response
type JawabanFileResponse struct {
	ID        int       `json:"id"`
	JawabanID int       `json:"jawaban_id"` // Wajib ada jawaban
	FileURL   string    `json:"file_url"`   // Wajib ada file URL
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
