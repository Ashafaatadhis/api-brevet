package validation

// PostManageGuru struct untuk validasi
type PostManageGuru struct {
	ID       uint   `json:"id"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Username string `json:"username" validate:"required,min=3,max=100"`
	RoleID   uint   `json:"role_id"`
	Nohp     string `json:"nohp" validate:"required,min=3,max=16"`
	Avatar   string `json:"avatar" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=255"`
}

// TableName untuk representasi ke table db
func (PostManageGuru) TableName() string {
	return "users"
}

// UpdateManageGuru struct untuk validasi, tanpa password
type UpdateManageGuru struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Username string `json:"username" validate:"required,min=3,max=100"`
	Nohp     string `json:"nohp" validate:"required,min=3,max=16"`
	Email    string `json:"email" validate:"required,email"`
}
