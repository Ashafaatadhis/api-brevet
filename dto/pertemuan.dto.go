package dto

import (
	"time"
)

// CreatePertemuanRequest untuk request
type CreatePertemuanRequest struct {
	Name string `json:"name" validate:"required"`
}

// EditPertemuanRequest untuk request
type EditPertemuanRequest struct {
	Name *string `json:"name" validate:"omitempty"`
}

// PertemuanResponse struct untuk response khusus menangani data pertemuan
type PertemuanResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	GrBatchID int    `json:"group_batches_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Materis []MateriResponse `json:"materis"`
}
