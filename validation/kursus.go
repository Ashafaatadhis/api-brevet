package validation

import (
	"new-brevet-be/models"
	"time"
)

// PostKursus adalah representasi tabel kursus di database
type PostKursus struct {
	Judul            string    `form:"judul" validate:"required,min=3,max=255"`
	JenisID          int       `form:"jenis_id" validate:"required"`
	KelasID          int       `form:"kelas_id" validate:"required"`
	DeskripsiSingkat string    `form:"deskripsi_singkat" validate:"omitempty,max=255"`
	Deskripsi        string    `form:"deskripsi" validate:"required"`
	Pembelajaran     string    `form:"pembelajaran" validate:"required"`
	Diperoleh        string    `form:"diperoleh" validate:"required"`
	CategoryID       int       `form:"category_id" validate:"required"`
	ThumbnailKursus  string    `form:"thumbnail_kursus" validate:"omitempty,url"`
	StartDate        time.Time `form:"start_date" validate:"required"`
	EndDate          time.Time `form:"end_date" validate:"required,gtefield=StartDate"`
	StartTime        time.Time `form:"start_time" validate:"required"`
	EndTime          time.Time `form:"end_time" validate:"required,gtfield=StartTime"`

	HariID []uint        `validate:"required,dive,required" form:"hari_id"`
	Hari   []models.Hari `form:"hari"`
}

// TableName untuk representasi ke table db
func (PostKursus) TableName() string {
	return "kursus"
}
