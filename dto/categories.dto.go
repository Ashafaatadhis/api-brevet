package dto

import "time"

// CategoryResponse struct untuk response khusus menangani data categories
type CategoryResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
