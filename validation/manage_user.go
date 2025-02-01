package validation

import "time"

// PostManageUser struct untuk validasi
type PostManageUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Username string `json:"username" validate:"required,min=3,max=100"`
	RoleID   int    `json:"role_id" validate:"required,oneof=1 2 3 4"`
	Nohp     string `json:"nohp" validate:"required,min=3,max=16"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=255"`

	// profile purpose
	Institusi string    `json:"institusi" validate:"required,max=100"`
	Asal      string    `json:"asal" validate:"required,max=100"`
	TglLahir  time.Time `json:"tgl_lahir" validate:"required"`
	Alamat    string    `json:"alamat" validate:"required,max=255"`
}

// TableName untuk representasi ke table db
func (PostManageUser) TableName() string {
	return "users"
}

// UpdateManageUser struct untuk validasi, tanpa password
type UpdateManageUser struct {
	Name     *string `json:"name" validate:"omitempty,min=3,max=100"`
	Username *string `json:"username" validate:"omitempty,min=3,max=100"`
	Nohp     *string `json:"nohp" validate:"omitempty,min=3,max=16"`
	Email    *string `json:"email" validate:"omitempty,email"`
	RoleID   *int    `json:"role_id" validate:"omitempty,oneof=1 2 3 4"`

	// profile purpose
	Institusi *string    `json:"institusi" validate:"omitempty,max=100"`
	Asal      *string    `json:"asal" validate:"omitempty,max=100"`
	TglLahir  *time.Time `json:"tgl_lahir" validate:"omitempty"`
	Alamat    *string    `json:"alamat" validate:"omitempty,max=255"`
}
