package dto

import "time"

// MateriResponse dto response
type MateriResponse struct {
	ID          int       `json:"id"`
	PertemuanID int       `json:"pertemuan_id"` // Relasi ke pertemuan
	Title       string    `json:"title"`
	Content     string    `json:"content"`  // Bisa berupa teks
	FileURL     string    `json:"file_url"` // Jika ada file (PDF, video, dsb)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateMateriRequest untuk request
type CreateMateriRequest struct {
	Title   string  `form:"title" validate:"required"`
	Content *string `form:"content" validate:"omitempty"`
}

// UpdateMateriRequest untuk request
type UpdateMateriRequest struct {
	Title   *string `form:"title" validate:"omitempty"`
	Content *string `form:"content" validate:"omitempty"`
}
