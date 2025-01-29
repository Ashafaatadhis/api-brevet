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
	RoleID    int       `json:"role_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi ke model Role
	Profile responseProfile `json:"profile"`
	Role    responseRole    `json:"role"`
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

type responseProfile struct {
	ID         int       `json:"id"`
	GolonganID *int      `json:"golongan_id"`
	UserID     *int      `json:"user_id"`
	NIM        *string   `json:"nim"`
	BuktiNIM   *string   `json:"bukti_nim"`
	NIK        *string   `json:"nik"`
	BuktiNIK   *string   `json:"bukti_nik"`
	Institusi  string    `json:"institusi"`
	Asal       string    `json:"asal"`
	TglLahir   time.Time `json:"tgl_lahir"`
	Alamat     string    `json:"alamat"`

	Golongan responseKategoriGolongan `json:"kategori_golongan"`
}

// TableName untuk representasi ke table profiles
func (responseProfile) TableName() string {
	return "profiles"
}

type responseKategoriGolongan struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TableName untuk representasi ke table kategori_golongans
func (responseKategoriGolongan) TableName() string {
	return "kategori_golongans"
}
