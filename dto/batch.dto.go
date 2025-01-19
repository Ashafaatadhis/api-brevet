package dto

import (
	"time"
)

// BatchResponse struct untuk response khusus menangani data batch
type BatchResponse struct {
	ID         int       `json:"id"`
	Judul      string    `validate:"required,min=3,max=100" json:"judul"`            // Judul wajib diisi, minimal 3 karakter, maksimal 100
	BukaBatch  time.Time `validate:"required" json:"buka_batch"`                     // Tanggal buka batch wajib diisi
	TutupBatch time.Time `validate:"required,gtefield=BukaBatch" json:"tutup_batch"` // Tanggal tutup batch wajib >= BukaBatch

	GroupBatches []GroupBatchResponse `json:"group_batches"`
}

// BatchRequest struct untuk request khusus menangani data batch
type BatchRequest struct {
	Judul      string    `validate:"required,min=3,max=100" json:"judul"`            // Judul wajib diisi, minimal 3 karakter, maksimal 100
	BukaBatch  time.Time `validate:"required" json:"buka_batch"`                     // Tanggal buka batch wajib diisi
	TutupBatch time.Time `validate:"required,gtefield=BukaBatch" json:"tutup_batch"` // Tanggal tutup batch wajib >= BukaBatch
}

// TableName untuk representasi ke table users
func (BatchResponse) TableName() string {
	return "batches"
}

// TableName untuk representasi ke table users
func (BatchRequest) TableName() string {
	return "batches"
}
