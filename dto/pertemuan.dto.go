package dto

import (
	"time"
)

// CreatePertemuanRequest untuk request
type CreatePertemuanRequest struct {
	GrBatchID int    `json:"group_batches_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
}

// EditPertemuanRequest untuk request
type EditPertemuanRequest struct {
	ID        int     `json:"-"` // ID dari params
	GrBatchID *int    `json:"group_batches_id" validate:"omitempty"`
	Name      *string `json:"name" validate:"omitempty"`
}

// SetID Method untuk di-set di middleware
func (r *EditPertemuanRequest) SetID(id int) {
	r.ID = id
}

// PertemuanResponse struct untuk response khusus menangani data pertemuan
type PertemuanResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	GrBatchID int    `json:"group_batches_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	GroupBatch GroupBatchResponse `json:"group_batches"`
}
