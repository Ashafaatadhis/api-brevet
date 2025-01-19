package models

// KategoriGolongan adalah representasi tabel kategory_golongan di database
type KategoriGolongan struct {
	ID   int    `gorm:"primaryKey;type:int unsigned;"`
	Name string `gorm:"size:100;not null;uniqueIndex"`
}
