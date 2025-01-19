package models

// KelasKursus adalah representasi tabel kelas_kursus di database
type KelasKursus struct {
	ID   int    `gorm:"primaryKey;type:int unsigned;"`
	Name string `gorm:"size:100;not null;uniqueIndex"`
}
