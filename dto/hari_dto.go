package dto

import "time"

// HariResponse struct untuk response khusus menangani data kuharirsus
type HariResponse struct {
	ID        int              `json:"id"`
	Nama      string           `json:"nama"`
	Kursus    []KursusResponse `json:"kursus,omitempty"` // Data kursus dalam bentuk ringkas
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// TableName untuk representasi ke table users
func (HariResponse) TableName() string {
	return "haris"
}
