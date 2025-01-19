package dto

import (
	"time"
)

// GroupBatchResponse struct untuk response khusus menangani data batchgroup
type GroupBatchResponse struct {
	ID        int  `json:"id"`
	TeacherID *int `json:"teacher_id"`
	BatchID   *int `json:"batch_id"`
	KursusID  *int `json:"kursus_id"`

	Teacher *responseUser   `json:"teacher"`
	Batch   *BatchResponse  `json:"batches"`
	Kursus  *responseKursus `json:"kursus"`
}

type responseUser struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Nohp      string `json:"nohp"`
	Avatar    string `json:"avatar"`
	RoleID    int    `json:"roleId"`
	Email     string `json:"email"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type responseKursus struct {
	ID               int       `json:"id"`
	TeacherID        *string   `json:"teacher_id"`
	Judul            string    `json:"judul"`
	JenisID          int       `json:"jenis_id"`
	KelasID          int       `json:"kelas_id"`
	DeskripsiSingkat string    `json:"deskripsi_singkat"`
	Deskripsi        string    `json:"deskripsi"`
	Pembelajaran     string    `json:"pembelajaran"`
	Diperoleh        string    `json:"diperoleh"`
	CategoryID       int       `json:"category_id"`
	ThumbnailKursus  string    `json:"thumbnail_kursus"`
	ThumbnailURL     string    `json:"thumbnail_url"`
	HargaAsli        float64   `json:"harga_asli"`
	HargaDiskon      float64   `json:"harga_diskon"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName untuk representasi ke table users
func (GroupBatchResponse) TableName() string {
	return "group_batches"
}

// TableName untuk representasi ke table users
func (responseKursus) TableName() string {
	return "kursus"
}

// TableName untuk representasi ke table users
func (responseUser) TableName() string {
	return "users"
}
