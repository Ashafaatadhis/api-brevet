package validation

// UserSetting struct untuk validasi, tanpa password
type UserSetting struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Username string `json:"username" validate:"required,min=3,max=100"`
	Nohp     string `json:"nohp" validate:"required,min=3,max=16"`
	Avatar   string `json:"avatar" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
}

// ChangePassword struct validasi rubah password
type ChangePassword struct {
	OldPassword string `json:"old_password" validate:"required,min=3,max=255"`
	NewPassword string `json:"new_password" validate:"required,min=3,max=255"`
}
