package dto

import (
	"new-brevet-be/models"
	"time"
)

// KursusResponse struct untuk response khusus menangani data kursus
type KursusResponse struct {
	ID int `json:"id"`

	Judul string `json:"judul"`

	DeskripsiSingkat string    `json:"deskripsi_singkat"`
	Deskripsi        string    `json:"deskripsi"`
	Pembelajaran     string    `json:"pembelajaran"`
	Diperoleh        string    `json:"diperoleh"`
	CategoryID       int       `json:"category_id"`
	ThumbnailKursus  string    `json:"thumbnail_kursus"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	StartTime        string    `json:"start_time"`
	EndTime          string    `json:"end_time"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Category models.Category `json:"category,omitempty"`

	GroupBatches []*GroupBatchResponse `json:"group_batches"`

	Hari []models.Hari `json:"hari"` // ID hari yang terkait
}

// TableName untuk representasi ke table users
func (KursusResponse) TableName() string {
	return "kursus"
}
