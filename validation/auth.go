package validation

// UserLogin struct type untuk body login
type UserLogin struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=3,max=255"`
}

// TableName untuk representasi ke table db
func (UserLogin) TableName() string {
	return "users"
}

// UserRegister struct type untuk body register
type UserRegister struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Username string `json:"username" validate:"required,min=3,max=100"`
	RoleID   int    `json:"role_id" validate:"required"`
	Nohp     string `json:"nohp" validate:"required,min=3,max=16"`
	Avatar   string `json:"avatar" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=255"`
}

// TableName untuk representasi ke table db
func (UserRegister) TableName() string {
	return "users"
}
