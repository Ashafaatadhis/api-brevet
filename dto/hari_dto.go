package dto

import "time"

// HariResponse struct untuk response khusus menangani data kuharirsus
type HariResponse struct {
	ID        int       `json:"id"`
	Nama      string    `json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
