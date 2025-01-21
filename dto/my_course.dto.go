package dto

import "time"

// MyCourseResponse struct untuk response
type MyCourseResponse struct {
	ID            int       `json:"id"`
	GrBatchID     int       `json:"group_batches_id"`
	JenisKursusID int       `json:"jenis_kursus_id"`
	URLConfirm    *string   `json:"url_confirm"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	GroupBatches GroupBatchResponse `json:"group_batches"`
}
