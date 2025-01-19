package handlers

import (
	"log"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/services"
	"new-brevet-be/utils"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAllBuyKursus adalah handler untuk route buy-kursus
func GetAllBuyKursus(c *fiber.Ctx) error {
	db := config.DB

	var purchase []models.Purchase

	var responsePurchase []dto.BuykursusResponse

	if err := db.Preload("JenisKursus").
		Preload("GroupBatches").
		Preload("GroupBatches.Teacher").
		Preload("GroupBatches.Batch").
		Preload("GroupBatches.Kursus").
		Preload("User").
		Preload("User.Role").
		Preload("StatusPayment").
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, harga")
		}).
		Find(&purchase).Error; err != nil {
		log.Println("Purchase not exist:", err)
		return utils.Response(c, fiber.StatusBadRequest, "Purchase not exist", nil, nil, nil)
	}

	// Automapping
	if err := dto_mapper.Map(&responsePurchase, purchase); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map purchasing response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "success retrieved all purchasing", responsePurchase, nil, nil)
}

// GetBuyKursus adalah handler untuk route buy-kursus/:id
func GetBuyKursus(c *fiber.Ctx) error {
	db := config.DB
	ID := c.Params("id")
	var purchase models.Purchase

	var responsePurchase dto.BuykursusResponse

	if err := db.Preload("JenisKursus").
		Preload("GroupBatches").
		Preload("GroupBatches.Teacher").
		Preload("GroupBatches.Batch").
		Preload("GroupBatches.Kursus").
		Preload("User").
		Preload("User.Role").
		Preload("StatusPayment").
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, harga")
		}).
		First(&purchase, ID).Error; err != nil {
		log.Println("Purchase not exist:", err)
		return utils.Response(c, fiber.StatusBadRequest, "Purchase not exist", nil, nil, nil)
	}

	// Automapping
	if err := dto_mapper.Map(&responsePurchase, purchase); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map purchasing response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "success retrieved all purchasing", responsePurchase, nil, nil)
}

// CreateBuyKursus adalah handler untuk route buy-kursus
func CreateBuyKursus(c *fiber.Ctx) error {

	db := config.DB
	body := c.Locals("body").(dto.BuyKursusRequest)
	token := c.Locals("user").(middlewares.User)

	tx := db.Begin()

	var user models.User

	if err := tx.Preload("Profile", func(db *gorm.DB) *gorm.DB {
		return db.Select("golongan_id, user_id")
	}).First(&user, token.ID).Error; err != nil {
		tx.Rollback()
		log.Println("User not exist:", err)
		return utils.Response(c, fiber.StatusBadRequest, "User not exist", nil, nil, nil)
	}

	confirmationCode, err := services.GenerateURLConfirm()
	if err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create Registration", nil, nil, nil)
	}

	var price models.Price
	if err := tx.Select("id, harga").Where("golongan_id = ?", &user.Profile.GolonganID).First(&price).Error; err != nil {
		tx.Rollback()
		log.Println("Cannot fetch price:", err)
		return utils.Response(c, fiber.StatusBadRequest, "Cannot fetch price", nil, nil, nil)
	}

	purchase := models.Purchase{
		GrBatchID:       body.GroupBatchesID,
		StatusPaymentID: 1, //pending
		PriceID:         price.ID,
		JenisKursusID:   body.JenisKursusID,
		UserID:          &token.ID,
		URLConfirm:      &confirmationCode,
	}

	if err := tx.Create(&purchase).Error; err != nil {
		tx.Rollback()
		log.Println("Failed to buy Kursus:", err)
		return utils.Response(c, fiber.StatusBadRequest, "Failed to buy Kursus", nil, nil, nil)
	}

	// Kirim kode pembayaran ke email
	if err := services.SendEmailCodePayment(user.Name, user.Email, purchase.URLConfirm); err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to send email", nil, nil, nil)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Failed to commit transaction:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to buy kursus", nil, nil, nil)
	}

	var responsePurchase dto.BuykursusResponse

	if err := db.Preload("JenisKursus").
		Preload("GroupBatches").
		Preload("GroupBatches.Teacher").
		Preload("GroupBatches.Batch").
		Preload("GroupBatches.Kursus").
		Preload("User").
		Preload("User.Role").
		Preload("StatusPayment").
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, harga")
		}).
		First(&purchase, purchase.ID).Error; err != nil {
		log.Println("Purchase not exist:", err)
		return utils.Response(c, fiber.StatusBadRequest, "Purchase not exist", nil, nil, nil)
	}

	// Automapping
	if err := dto_mapper.Map(&responsePurchase, purchase); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Buy Kursus buy successfully", responsePurchase, nil, nil)

}

// EditBuyKursus adalah handler untuk route buy-kursus/:id (only status payment can edit)
func EditBuyKursus(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(dto.EditBuyKursus)
	ID := c.Params("id")

	tx := db.Begin()

	var purchase models.Purchase

	if err := tx.Preload("User").
		First(&purchase, ID).Error; err != nil {
		tx.Rollback()
		log.Println("Purchase not exist:", err)
		return utils.Response(c, fiber.StatusBadRequest, "Purchase not exist", nil, nil, nil)
	}

	purchase.StatusPaymentID = body.StatusPaymentID

	if err := tx.Save(&purchase).Error; err != nil {
		tx.Rollback()
		log.Println("Failed to change status purchase:", err)
		return utils.Response(c, fiber.StatusBadRequest, "Failed to change status purchase", nil, nil, nil)
	}

	if body.StatusPaymentID == 2 { // kalau lunas kirim email
		// Kirim kode pembayaran ke email
		if err := services.SendEmailConfirmAccount(purchase.User.Username, purchase.User.Email); err != nil {
			tx.Rollback()
			log.Println("Failed to send email:", err)
			return utils.Response(c, fiber.StatusInternalServerError, "Failed to send email", nil, nil, nil)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Failed to commit transaction:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to buy kursus", nil, nil, nil)
	}

	var responsePurchase dto.BuykursusResponse

	if err := db.Preload("JenisKursus").
		Preload("GroupBatches").
		Preload("GroupBatches.Teacher").
		Preload("GroupBatches.Batch").
		Preload("GroupBatches.Kursus").
		Preload("User").
		Preload("User.Role").
		Preload("StatusPayment").
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, harga")
		}).
		First(&purchase, purchase.ID).Error; err != nil {
		log.Println("Purchase not exist:", err)
		return utils.Response(c, fiber.StatusBadRequest, "Purchase not exist", nil, nil, nil)
	}

	// Automapping
	if err := dto_mapper.Map(&responsePurchase, purchase); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Edit Buy Kursus Successfully", responsePurchase, nil, nil)
}

// ConfirmPayment adalah handler untuk route confirm-payement/:id
func ConfirmPayment(c *fiber.Ctx) error {
	db := config.DB

	kodePembayaran := c.Params("id")

	var purchase models.Purchase

	// upload gambar
	path := "bukti"

	// Cari data purchase berdasarkan ID
	if err := db.Where("url_confirm = ?", kodePembayaran).First(&purchase).Error; err != nil {
		return utils.Response(c, fiber.StatusNotFound, "Kode pembayaran not found", nil, nil, nil)
	}

	if purchase.BuktiBayar != nil {
		return utils.Response(c, fiber.StatusNotFound, "sudah upload bukti bayar", nil, nil, nil)
	}

	// Proses upload file
	dataBayar, err := utils.UploadFile(c, "bukti_bayar", path)
	if err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "Failed to upload Bukti Bayar", nil, nil, nil)
	} else if dataBayar != nil {
		purchase.BuktiBayar = dataBayar
	}

	if err := db.Model(&purchase).Update("bukti_bayar", *purchase.BuktiBayar).Error; err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "Failed to update Pembayaran", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Sukses upload bukti bayar", nil, nil, nil)
}
