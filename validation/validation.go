package validation

import (
	"fmt"
	"new-brevet-be/config"
	"new-brevet-be/middlewares"
	"new-brevet-be/utils"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Validate adalah fungsi utama untuk memvalidasi request body
func Validate[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload T

		// Parse request body
		if err := c.BodyParser(&payload); err != nil {
			logrus.Error("ERROR: fail to parsing request body:", err.Error())
			return utils.Response(c, fiber.StatusBadRequest, err.Error(), nil, nil, nil)
		}

		idParam := c.Params("id")
		// userData := c.Locals("user").(middlewares.User)
		userData, ok := c.Locals("user").(middlewares.User)
		if !ok {
			logrus.Warn("WARNING: c.locals user data not found")
			userData = middlewares.User{} // Set default userData jika kosong
		}

		// Default ID menggunakan userData.ID jika idParam kosong
		var finalID int
		var err error

		if idParam != "" {
			finalID, err = strconv.Atoi(idParam)
			if err != nil {
				logrus.Error("ERROR: error parsing idparam to integer:", err.Error())
				return utils.Response(c, fiber.StatusBadRequest, "Invalid ID parameter", nil, nil, nil)
			}
		} else {
			finalID = userData.ID // Gunakan user ID jika idParam kosong
		}

		// Cek apakah payload memiliki method SetID, jika iya, set dengan finalID
		if idSetter, ok := any(&payload).(interface{ SetID(int) }); ok {
			idSetter.SetID(finalID)
		}

		// Validasi menggunakan validator
		// Dapatkan validator global
		validate := InitValidator()

		if idParam != "" {

			validate.RegisterValidation("unique_except", ValidateUniqueExcept(idParam))
		}

		if err := validate.Struct(&payload); err != nil {
			// Jika ada error validasi, kirimkan response error dengan detailnya
			var validationErrors []string
			for _, err := range err.(validator.ValidationErrors) {

				if err.Tag() == "exists" {
					// Menampilkan error dengan format kustom

					validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' is invalid. id %d does not exist in the database", err.Field(), err.Value()))
				} else if err.Tag() == "unique" || err.Tag() == "unique_except" {
					// Menampilkan error dengan format kustom

					validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' is invalid. id %d already in the database", err.Field(), err.Value()))
				} else {

					validationErrors = append(validationErrors, err.Error())
				}

			}
			logrus.Error("ERROR: validation failed:", err.Error())
			return utils.Response(c, fiber.StatusBadRequest, "Validation failed", validationErrors, nil, nil)
		}

		// Simpan payload yang sudah divalidasi ke dalam `Locals` dengan key `body`
		c.Locals("body", payload)
		logrus.Info("Success validation")
		// Lanjutkan ke handler berikutnya jika tidak ada error
		return c.Next()
	}
}

// ValidateExists function for checking if a record exists in a general table and column
func ValidateExists(fl validator.FieldLevel) bool {

	db := config.DB
	// Extract table and column information from the tag
	param := fl.Param()
	parts := strings.Split(param, ".")
	if len(parts) != 2 {
		// Tag format must be "table.column"
		logrus.Warn("WARNING: Invalid validation parameter format, expected 'table.column'")

		return false
	}

	table, column := parts[0], parts[1]

	// Get the value of the field being validated
	value := fl.Field().Interface()

	// Check if the record exists in the database
	var count int64
	query := fmt.Sprintf("%s = ?", column)
	if err := db.Table(table).Where(query, value).Count(&count).Error; err != nil {
		logrus.Warn("ERROR: Error querying database:", err.Error())
		return false
	}

	logrus.Info("Success validate exists")
	// Return true if the record exists
	return count > 0

}

// ValidateExistsExcept function for checking if a record exists with an exception
func ValidateExistsExcept(idNow string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		db := config.DB
		param := fl.Param() // Contoh: "table.column"
		parts := strings.Split(param, ".")
		if len(parts) != 2 {
			logrus.Warn("WARNING: Invalid validation parameter format, expected 'table.column'")
			return false
		}

		table := parts[0]
		column := parts[1]
		fieldValue := fl.Field().Interface()

		// Cek apakah field ada di database, kecuali ID saat ini
		var count int64
		query := fmt.Sprintf("%s = ? AND id != ?", column)
		if err := db.Table(table).Where(query, fieldValue, idNow).Count(&count).Error; err != nil {
			logrus.Warn("ERROR: Error querying database:", err.Error())
			return false
		}

		logrus.Info("Success validate exists except")
		return count == 0
	}

}

// ValidateUnique function for checking if a record was exist with an exception
func ValidateUnique(fl validator.FieldLevel) bool {
	db := config.DB
	param := fl.Param() // Format param: "table.column"
	parts := strings.Split(param, ".")
	if len(parts) != 2 {
		logrus.Warn("WARNING: Invalid validation parameter format, expected 'table.column'")
		return false
	}

	table := parts[0]
	column := parts[1]
	fieldValue := fl.Field().Interface()

	// Cek apakah field sudah ada di database
	var count int64
	if err := db.Table(table).Where(fmt.Sprintf("%s = ?", column), fieldValue).Count(&count).Error; err != nil {
		logrus.Warn("ERROR: Error querying database:", err.Error())
		return false
	}
	logrus.Info("Success validate unique")
	return count == 0

}

// ValidateUniqueExcept function for checking if a record was exist except id given with an exception
func ValidateUniqueExcept(idNow string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		db := config.DB
		param := fl.Param() // Format param: "table.column"
		parts := strings.Split(param, ".")
		if len(parts) != 2 {
			logrus.Warn("WARNING: Invalid validation parameter format, expected 'table.column'")
			return false
		}

		table := parts[0]
		column := parts[1]
		fieldValue := fl.Field().Interface()

		// Cek apakah field sudah ada di database, kecuali untuk ID tertentu
		var count int64
		query := fmt.Sprintf("%s = ? AND id != ?", column)
		if err := db.Table(table).Where(query, fieldValue, idNow).Count(&count).Error; err != nil {
			logrus.Warn("ERROR: Error querying database:", err.Error())
			return false
		}

		logrus.Info("Success validate unique except")
		return count == 0
	}
}
