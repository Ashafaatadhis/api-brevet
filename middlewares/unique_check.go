package middlewares

import (
	"new-brevet-be/config"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// UserUniqueCheck adalah middleware untuk memastikan bahwa username, nohp, dan email unik
func UserUniqueCheck[T any](c *fiber.Ctx) error {

	db := config.DB
	body := c.Locals("body").(T) // Ambil payload dari Locals
	userInterface := c.Locals("user")

	// Gunakan reflect untuk membaca field
	v := reflect.ValueOf(body)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Ambil nilai field yang akan divalidasi
	username := v.FieldByName("Username").String()
	nohp := v.FieldByName("Nohp").String()
	email := v.FieldByName("Email").String()

	// Query untuk memeriksa duplikasi
	var existingUser models.User
	query := db.Where("username = ? OR nohp = ? OR email = ?", username, nohp, email)

	userID := c.Params("id")
	if userID != "" {
		query = query.Not("id = ?", userID)
		userInterface = nil
	}

	// Tambahkan pengecualian untuk ID jika disediakan
	if userInterface != nil && c.Method() != "POST" {
		user, ok := userInterface.(User)
		if ok {

			query = query.Not("id = ?", user.ID)
		}

	}

	// Eksekusi query
	if err := query.First(&existingUser).Error; err == nil {
		var conflictField string

		switch {
		case existingUser.Username == username:
			conflictField = "Username is already taken"
		case existingUser.Nohp == nohp:
			conflictField = "Phone number is already registered"
		case existingUser.Email == email:
			conflictField = "Email is already registered"
		}

		return utils.Response(c, fiber.StatusBadRequest, conflictField, nil, nil, nil)
	}

	// Lanjutkan ke handler berikutnya jika tidak ada konflik
	return c.Next()

}

// KursusUniqueCheck adalah middleware untuk memastikan bahwa Judul kursus unik
func KursusUniqueCheck[T any](c *fiber.Ctx) error {

	db := config.DB
	body := c.Locals("body").(T) // Ambil payload dari Locals

	// Gunakan reflect untuk membaca field
	v := reflect.ValueOf(body)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Ambil nilai field yang akan divalidasi (Judul)
	judul := v.FieldByName("Judul").String()

	// Query untuk memeriksa duplikasi Judul
	var existingKursus models.Kursus
	query := db.Where("judul = ?", judul)

	// Ambil user ID dari params jika ada (misalnya untuk PUT/PATCH)
	kursusID := c.Params("id")
	if kursusID != "" {
		query = query.Not("id = ?", kursusID)
	}

	// Eksekusi query
	if err := query.First(&existingKursus).Error; err == nil {
		// Jika Judul sudah ada
		return utils.Response(c, fiber.StatusBadRequest, "Judul kursus sudah terdaftar", nil, nil, nil)
	}

	// Lanjutkan ke handler berikutnya jika tidak ada konflik
	return c.Next()
}

// BatchUniqueCheck adalah middleware untuk memastikan bahwa Judul batch unik
func BatchUniqueCheck[T any](c *fiber.Ctx) error {

	db := config.DB
	body := c.Locals("body").(T) // Ambil payload dari Locals

	// Gunakan reflect untuk membaca field
	v := reflect.ValueOf(body)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Ambil nilai field yang akan divalidasi (Judul)
	judul := v.FieldByName("Judul").String()

	// Query untuk memeriksa duplikasi Judul
	var existingBatch models.Batch
	query := db.Where("judul = ?", judul)

	// Ambil user ID dari params jika ada (misalnya untuk PUT/PATCH)
	batchID := c.Params("id")
	if batchID != "" {
		query = query.Not("id = ?", batchID)
	}

	// Eksekusi query
	if err := query.First(&existingBatch).Error; err == nil {
		// Jika Judul sudah ada
		return utils.Response(c, fiber.StatusBadRequest, "Judul batch sudah terdaftar", nil, nil, nil)
	}

	// Lanjutkan ke handler berikutnya jika tidak ada konflik
	return c.Next()
}

// BuyBatchUniqueCheck adalah middleware untuk memastikan bahwa kursus belum ada di transaksi
func BuyBatchUniqueCheck[T any](c *fiber.Ctx) error {

	db := config.DB
	body := c.Locals("body").(T) // Ambil payload dari Locals

	// Gunakan reflect untuk membaca field
	v := reflect.ValueOf(body)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Ambil nilai field yang akan divalidasi (Judul)
	idGroupBatch := v.FieldByName("GroupBatchesID").Int()

	// Query untuk memeriksa duplikasi Judul
	var existingGroupBatch models.Purchase
	query := db.Where("gr_batch_id = ?", idGroupBatch)

	// Eksekusi query
	if err := query.First(&existingGroupBatch).Error; err == nil {
		// Jika Judul sudah ada
		if existingGroupBatch.StatusPaymentID == 2 {
			return utils.Response(c, fiber.StatusBadRequest, "Kursus Sudah Dibeli", nil, nil, nil)
		} else if existingGroupBatch.StatusPaymentID == 1 {
			return utils.Response(c, fiber.StatusBadRequest, "Kursus Sudah Dibeli, segera lakukan pembayaran", nil, nil, nil)
		}
	}

	// Lanjutkan ke handler berikutnya jika tidak ada konflik
	return c.Next()
}

// // RegistrationUniqueCheck adalah middleware untuk memastikan bahwa email dan nohp saat registrasi unik
// func RegistrationUniqueCheck[T any](c *fiber.Ctx) error {

// 	db := config.DB
// 	body := c.Locals("body").(T) // Ambil payload dari Locals

// 	// Gunakan reflect untuk membaca field
// 	v := reflect.ValueOf(body)
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}

// 	// Ambil nilai field yang akan divalidasi (Judul)
// 	email := v.FieldByName("Email").String()
// 	nohp := v.FieldByName("Nohp").String()

// 	// Query untuk memeriksa duplikasi Judul
// 	var existingBatch models.User
// 	query := db.Where("email = ? OR no_hp = ?", email, nohp)

// 	// Eksekusi query
// 	if err := query.First(&existingBatch).Error; err == nil {
// 		// Jika Judul sudah ada
// 		return utils.Response(c, fiber.StatusBadRequest, "email atau nohp sudah terdaftar", nil, nil, nil)
// 	}

// 	var existingUser models.User
// 	queryUser := db.Where("email = ? OR nohp = ?", email, nohp)

// 	// Eksekusi query
// 	if err := query.First(&existingBatch).Error; err == nil {
// 		// Jika Judul sudah ada
// 		return utils.Response(c, fiber.StatusBadRequest, "email atau nohp sudah terdaftar", nil, nil, nil)
// 	}

// 	// Eksekusi query
// 	if err := queryUser.First(&existingUser).Error; err == nil {
// 		// Jika Judul sudah ada
// 		return utils.Response(c, fiber.StatusBadRequest, "email atau nohp sudah terdaftar", nil, nil, nil)
// 	}

// 	// Lanjutkan ke handler berikutnya jika tidak ada konflik
// 	return c.Next()
// }
