package validation

import "time"

// UserSetting struct untuk validasi, tanpa password
type UserSetting struct {
	Name     string `form:"name" validate:"required,min=3,max=100"`
	Username string `form:"username" validate:"required,min=3,max=100"`
	Nohp     string `form:"nohp" validate:"required,min=3,max=16"`
	Avatar   string `form:"avatar" validate:"omitempty,min=3,max=100"`
	Email    string `form:"email" validate:"required,email"`

	// profile purpose
	Institusi string    `form:"institusi" validate:"required,max=100"`
	Asal      string    `form:"asal" validate:"required,max=100"`
	TglLahir  time.Time `form:"tgl_lahir" validate:"required"`
	Alamat    string    `form:"alamat" validate:"required,max=255"`
}

// ChangePassword struct validasi rubah password
type ChangePassword struct {
	OldPassword string `json:"old_password" validate:"required,min=3,max=255"`
	NewPassword string `json:"new_password" validate:"required,min=3,max=255"`
}
