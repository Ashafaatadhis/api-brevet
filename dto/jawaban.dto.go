package dto

import (
	"mime/multipart"
	"time"
)

// SubmitJawabanRequest untuk request
type SubmitJawabanRequest struct {
	Answer *string                 `form:"answer" validate:"omitempty"` // Bisa kosong kalau hanya upload file
	Files  []*multipart.FileHeader `form:"files" validate:"omitempty"`  // Bisa lebih dari satu file
}

// JawabanResponse untuk response
type JawabanResponse struct {
	ID        int       `json:"id"`
	TugasID   int       `json:"tugas_id"`
	UserID    int       `json:"user_id"` // One User, One Answer Per Tugas
	Answer    string    `json:"answer"`
	Score     *int      `json:"score,omitempty"`
	Feedback  *string   `json:"feedback,omitempty"`
	IsLate    bool      `json:"is_late"`
	IsGraded  bool      `json:"is_graded"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	JawabanFile []*JawabanFileResponse `json:"jawaban_file,omitempty"`
}
