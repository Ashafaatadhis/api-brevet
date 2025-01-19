package validation

// PostManageUser struct untuk validasi
type PostManageUser struct {
	ID       uint   `json:"id"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Username string `json:"username" validate:"required,min=3,max=100"`
	RoleID   uint   `json:"role_id" validate:"required,oneof=1 2 3 4"`
	Nohp     string `json:"nohp" validate:"required,min=3,max=16"`
	Avatar   string `json:"avatar" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=255"`
}

// TableName untuk representasi ke table db
func (PostManageUser) TableName() string {
	return "users"
}

// UpdateManageUser struct untuk validasi, tanpa password
type UpdateManageUser struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Username string `json:"username" validate:"required,min=3,max=100"`
	Nohp     string `json:"nohp" validate:"required,min=3,max=16"`
	Email    string `json:"email" validate:"required,email"`
}
