package utils

import "new-brevet-be/config"

// GetValidSortFields fungsi untuk mendapatkan column apa aja yang ada di table yang diberikan
func GetValidSortFields(model interface{}) (map[string]bool, error) {
	db := config.DB
	columns, err := db.Migrator().ColumnTypes(model)
	if err != nil {
		return nil, err
	}

	validSortFields := make(map[string]bool)
	for _, column := range columns {
		validSortFields[column.Name()] = true
	}

	return validSortFields, nil
}
