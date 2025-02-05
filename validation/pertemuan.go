package validation

import (
	"log"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/utils"

	"github.com/go-playground/validator/v10"
)

// CreatePertemuanUniqueValidation memeriksa apakah kombinasi GroupBatchesID dan Name sudah ada di database.
func CreatePertemuanUniqueValidation(sl validator.StructLevel) {
	// Ambil request yang sedang divalidasi
	req := sl.Current().Interface().(dto.CreatePertemuanRequest)
	db := config.DB

	// Cek apakah kombinasi group_batches_id dan name sudah ada di database
	var count int64
	err := db.Table("pertemuans").
		Where("gr_batch_id = ? AND name = ?", req.GrBatchID, req.Name).
		Count(&count).Error

	if err != nil {
		// Jika ada error pada query database, kita akan laporkan error generik
		sl.ReportError(req.Name, "Name", "name", "unique_check", "gagal memeriksa kombinasi group_batches_id dan name")
		return
	}

	if count > 0 {
		// Jika sudah ada, laporkan error pada field "Name"
		sl.ReportError(req.Name, "Name", "name", "unique", "kombinasi group_batches_id dan name sudah ada")
	}

}

// EditPertemuanUniqueValidation memeriksa apakah kombinasi GroupBatchesID dan Name sudah ada di database.
func EditPertemuanUniqueValidation(sl validator.StructLevel) {
	// Ambil request yang sedang divalidasi

	req := sl.Current().Interface().(dto.EditPertemuanRequest)
	db := config.DB

	// Cek apakah kombinasi group_batches_id dan name sudah ada di database
	var count int64
	err := db.Table("pertemuans").
		Where("id != ? AND gr_batch_id = ? AND name = ?", req.ID, utils.GetIntValue(req.GrBatchID), utils.GetStringValue(req.Name)).
		Count(&count).Error

	if err != nil {
		// Jika ada error pada query database, kita akan laporkan error generik
		sl.ReportError(req.Name, "Name", "name", "unique_check", "gagal memeriksa kombinasi group_batches_id dan name")
		return
	}

	if count > 0 {
		// Jika sudah ada, laporkan error pada field "Name"
		sl.ReportError(req.Name, "Name", "name", "unique", "kombinasi group_batches_id dan name sudah ada")
	}
	log.Print("kacau", err, count, req.ID)

}
