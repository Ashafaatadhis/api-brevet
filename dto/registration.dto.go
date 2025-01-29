package dto

import (
	"time"
)

// CreateRegistrationRequest struct untuk response khusus menangani request
type CreateRegistrationRequest struct {
	Name     string `form:"name" validate:"required,max=100"`
	Username string `form:"username" validate:"required,max=100"`
	Password string `form:"password" validate:"required,min=3,max=255"`
	Nohp     string `form:"nohp" form:"nohp" validate:"required,max=14"`
	Avatar   string `form:"avatar"`
	Email    string `form:"email" validate:"required,email,max=100"`

	NIM       *string `form:"nim" validate:"omitempty,max=15"`
	BuktiNIM  *string `form:"bukti_nim" validate:"omitempty,max=100"`
	NIK       *string `form:"nik" validate:"omitempty,max=20"`
	BuktiNIK  *string `form:"bukti_nik" validate:"omitempty,max=100"`
	Institusi string  `form:"institusi" validate:"required,max=100"`

	Asal     string    `form:"asal" validate:"required,max=100"`
	TglLahir time.Time `form:"tgl_lahir" validate:"required"`

	Alamat string `form:"alamat" validate:"required,max=255"`
	// URLConfirm *string   `form:"url_confirm" validate:"omitempty,max=255"`
	// BuktiBayar *string `form:"bukti_bayar" validate:"omitempty,max=100"`
	// Amount *string `form:"amount" validate:"omitempty,max=7"`
}

// EditRegistrationRequest struct untuk response khusus menangani request
type EditRegistrationRequest struct {
	GolonganID *int `json:"golongan_id" validate:"omitempty,exists=kategori_golongans.id"`
}

type profile struct {
	ID         int       `json:"id"`
	GolonganID *int      `json:"golongan_id"`
	UserID     *int      `json:"user_id"` // Pointer untuk mendukung nilai NULL
	NIM        *string   `json:"nim"`
	BuktiNIM   *string   `json:"bukti_nim"`
	NIK        *string   `json:"nik"`
	BuktiNIK   *string   `json:"bukti_nik"`
	Institusi  string    `json:"institusi"`
	Asal       string    `json:"asal"`
	TglLahir   time.Time `json:"tgl_lahir"`
	Alamat     string    `json:"alamat"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegistrationResponse struct untuk response khusus menangani response
type RegistrationResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Nohp      string    `json:"nohp"`
	Avatar    string    `json:"avatar"`
	Email     string    `json:"email"`
	Profile   profile   `json:"profile"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName untuk representasi ke table registrations
func (RegistrationResponse) TableName() string {
	return "registrations"
}

// TableName untuk representasi ke table registrations
func (CreateRegistrationRequest) TableName() string {
	return "registrations"
}

// TableName untuk representasi ke table registrations
func (EditRegistrationRequest) TableName() string {
	return "registrations"
}
