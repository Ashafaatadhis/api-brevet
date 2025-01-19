package models

// Category adalah representasi tabel categories di database
type Category struct {
	ID   int    `gorm:"primaryKey;type:int unsigned;"`
	Name string `gorm:"size:100;not null;uniqueIndex"`
}
