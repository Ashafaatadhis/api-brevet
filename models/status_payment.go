package models

// StatusPayment adalah representasi tabel status_payments di database
type StatusPayment struct {
	ID   int    `gorm:"primaryKey;type:int unsigned;" json:"id"`
	Name string `gorm:"size:100;not null;uniqueIndex" json:"name"`
}
