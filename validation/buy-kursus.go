package validation

import (
	"new-brevet-be/config"
	"new-brevet-be/dto"

	"github.com/go-playground/validator/v10"
)

// CreateBuyKursusUniqueValidation memeriksa apakah kombinasi GroupBatchesID sudah ada di database.
func CreateBuyKursusUniqueValidation(sl validator.StructLevel) {
	// Ambil request yang sedang divalidasi
	req := sl.Current().Interface().(dto.BuyKursusRequest)
	db := config.DB

	// Cek apakah kombinasi group_batches_id dan name sudah ada di database
	var count int64
	err := db.Table("purchases").
		Where("gr_batch_id = ? AND user_id = ?", req.GroupBatchesID, req.ID).
		Count(&count).Error

	if err != nil {
		// Jika ada error pada query database, kita akan laporkan error generik
		sl.ReportError(req.GroupBatchesID, "GroupBatchesID", "groupBatchesID", "unique_check", "gagal memeriksa kombinasi group_batches_id dan id")
		return
	}

	if count > 0 {
		// Jika sudah ada, laporkan error pada field "Name"
		sl.ReportError(req.GroupBatchesID, "GroupBatchesID", "groupBatchesID", "unique", "kombinasi group_batches_id dan id sudah ada")
	}

}
