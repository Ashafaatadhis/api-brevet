package dto

import (
	"time"
)

// ResponseUser struct untuk response khusus menangani data user
type ResponseUser struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Nohp      string    `json:"nohp"`
	Avatar    string    `json:"avatar"`
	RoleID    int       `json:"roleId"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi ke model Role
	Role responseRole `json:"role"`
}

// TableName untuk representasi ke table users
func (ResponseUser) TableName() string {
	return "users"
}

type responseRole struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TableName untuk representasi ke table roles
func (responseRole) TableName() string {
	return "roles"
}
