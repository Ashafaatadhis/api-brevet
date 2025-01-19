package models

// Price adalah representasi tabel Prices di database
type Price struct {
	ID         int `gorm:"primaryKey;type:int unsigned;" json:"id"`
	GolonganID int `gorm:"omitempty;type:int unsigned;" json:"golongan_id"`
	Harga      int `gorm:"not null" json:"harga"`

	KategoriGolongan KategoriGolongan `gorm:"foreignKey:GolonganID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
