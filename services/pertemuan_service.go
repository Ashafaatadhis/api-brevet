package services

import "new-brevet-be/config"

// IsPertemuanUnique Cek apakah kombinasi GroupBatchID dan Name sudah ada
func IsPertemuanUnique(grBatchID int, name string, excludeID *int) (bool, error) {
	db := config.DB
	var count int64

	query := db.Table("pertemuans").Where("gr_batch_id = ? AND name = ?", grBatchID, name)

	// Jika validasi untuk Edit, kita harus mengecualikan ID yang sedang diedit
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
