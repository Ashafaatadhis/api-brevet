package models

// JenisKursus adalah representasi tabel jenis_kursus di database
type JenisKursus struct {
	ID   int    `gorm:"primaryKey;type:int unsigned;" json:"id"`
	Name string `gorm:"size:100;not null;uniqueIndex" json:"name"`
}
