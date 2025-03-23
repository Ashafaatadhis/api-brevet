package dto

import (
	"mime/multipart"
	"time"
)

// TugasResponse dto response
type TugasResponse struct {
	ID          int       `json:"id"`
	PertemuanID int       `json:"pertemuan_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Type        string    `json:"type"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`

	Jawaban   []*JawabanResponse   `json:"jawaban,omitempty"`    // Bisa banyak file
	TugasFile []*TugasFileResponse `json:"tugas_file,omitempty"` // Bisa banyak file
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

// CreateTugasRequest untuk request
type CreateTugasRequest struct {
	Title       string                  `form:"title" validate:"required"`                 // Judul wajib diisi
	Description string                  `form:"description" validate:"required"`           // Deskripsi opsional
	Type        string                  `form:"type" validate:"required,oneof=essay file"` // Hanya bisa "essay" atau "file"
	StartAt     time.Time               `form:"start_at" validate:"required"`
	EndAt       time.Time               `form:"end_at" validate:"required,gtefield=StartAt"`
	Files       []*multipart.FileHeader `form:"files" validate:"omitempty"` // Bisa banyak file
}

// UpdateTugasRequest untuk request
type UpdateTugasRequest struct {
	Title       *string                 `form:"title" validate:"omitempty"`                 // Judul wajib diisi
	Description *string                 `form:"description" validate:"omitempty"`           // Deskripsi opsional
	Type        *string                 `form:"type" validate:"omitempty,oneof=essay file"` // Hanya bisa "essay" atau "file"
	StartAt     *time.Time              `form:"start_at" validate:"omitempty"`
	EndAt       *time.Time              `form:"end_at" validate:"omitempty,gtefield=StartAt"`
	Files       []*multipart.FileHeader `form:"files" validate:"omitempty"` // Bisa banyak file

}
