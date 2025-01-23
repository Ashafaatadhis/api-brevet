package dto

// MappingPengajarRequest struct request
type MappingPengajarRequest struct {
	TeacherID int `json:"teacher_id" validate:"required"`
}
