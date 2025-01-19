package validation

import "new-brevet-be/models"

// CreateRegistration struct untuk validation yang menangani route registration
type CreateRegistration struct {
	KursusID []uint          `json:"kursus_id" validate:"required,dive,required"`
	Kursus   []models.Kursus `json:"kursus,omitempty"` // Relasi many-to-many dengan Kursus
}
