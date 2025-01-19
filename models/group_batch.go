package models

// GroupBatch adalah representasi tabel group_batches di database
type GroupBatch struct {
	ID        int  `gorm:"primaryKey;autoIncrement;type:int unsigned;"`
	TeacherID *int `gorm:"omitempty;type:int unsigned;"`
	BatchID   *int `gorm:"omitempty;type:int unsigned;"`
	KursusID  *int `gorm:"omitempty;type:int unsigned;"`

	Teacher *User   `gorm:"foreignKey:TeacherID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Batch   *Batch  `gorm:"foreignKey:BatchID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Kursus  *Kursus `gorm:"foreignKey:KursusID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
