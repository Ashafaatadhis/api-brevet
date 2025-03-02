package handlers

import (
	"fmt"
	"log"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/services"
	"new-brevet-be/utils"
	"strings"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// GetAllBuyKursus adalah handler untuk route buy-kursus
func GetAllBuyKursus(c *fiber.Ctx) error {
	db := config.DB
	user := c.Locals("user").(middlewares.User)

	// Ambil query parameters
	search := c.Query("q", "")            // Pencarian (default kosong)
	sort := c.Query("sort", "id")         // Sorting field (default "id")
	order := c.Query("order", "asc")      // Urutan sorting (default "asc")
	selectFields := c.Query("select", "") // Field yang diinginkan (e.g., name, id)
	limit := c.QueryInt("limit", 10)      // Batas jumlah data (default 10)
	page := c.QueryInt("page", 1)         // Halaman (default 1)

	// Pagination offset
	offset := (page - 1) * limit

	// Ambil valid sort fields secara otomatis dari tabel
	validSortFields, err := utils.GetValidSortFields(&models.Purchase{})
	if err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get valid sort fields", nil, nil, err.Error())
	}

	// Validasi sort dan order
	if !validSortFields[sort] {
		sort = "id" // Default sorting field
	}
	if order != "asc" && order != "desc" {
		order = "asc" // Default order
	}

	var purchase []models.Purchase

	var responsePurchase []dto.BuykursusResponse

	query := db.Model(&models.Purchase{}).
		Joins("JOIN group_batches ON group_batches.id = purchases.gr_batch_id").
		Joins("JOIN kursus ON kursus.id = group_batches.kursus_id")

	// Apply search query
	if search != "" {
		query = query.Where("kursus.judul LIKE ?", "%"+search+"%")
	}

	query = query.Preload("GroupBatches").
		Preload("GroupBatches.Teacher").
		Preload("GroupBatches.Batch").
		Preload("GroupBatches.Kursus").
		Preload("User").
		Preload("User.Role").
		Preload("StatusPayment").
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, harga")
		})

		// logic ROLE
	if user.Role == "siswa" {

		query.Where("user_id = ?", user.ID)
	}

	// Apply select fields
	if selectFields != "" {
		// Pisahkan field berdasarkan koma (e.g., "name,id")
		fields := strings.Split(selectFields, ",")
		query = query.Select(fields)
	}

	// Apply sorting
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	// Hitung total data sebelum pagination
	var totalData int64
	if err := query.Count(&totalData).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to count total data", nil, nil, err.Error())
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&purchase).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get buy kursus", nil, nil, err.Error())
	}

	if err := copier.CopyWithOption(&responsePurchase, &purchase, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}

	// Set struct kosong menjadi nil
	utils.TransformResponse(&responsePurchase)

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "success retrieved all purchasing", responsePurchase, meta, nil)

}

// GetBuyKursus adalah handler untuk route buy-kursus/:id
func GetBuyKursus(c *fiber.Ctx) error {
	db := config.DB
	ID := c.Params("id")
	user := c.Locals("user").(middlewares.User)

	var purchase models.Purchase

	var responsePurchase dto.BuykursusResponse

	query := db

	if user.Role == "siswa" {
		query = query.Where("user_id = ?", user.ID)
	}
	query = query.Preload("GroupBatches").
		Preload("GroupBatches.Teacher").
		Preload("GroupBatches.Batch").
		Preload("GroupBatches.Kursus").
		Preload("User").
		Preload("User.Role").
		Preload("StatusPayment").
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, harga")
		})

	if err := query.First(&purchase, ID).Error; err != nil {
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

	// Cek kuota batch sebelum user membeli kursus
	hasQuota, err := services.CheckBatchQuota(tx, body.GroupBatchesID)
	if err != nil {
		tx.Rollback()

		log.Print(err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to check batch capacity", nil, nil, nil)
	}
	if !hasQuota {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Batch is full, cannot purchase", nil, nil, nil)
	}

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
		log.Println("Cannot fetch price, this happens cause this account gologang_id is nil:", err)
		return utils.Response(c, fiber.StatusBadRequest, "Your account non golongan", nil, nil, nil)
	}

	purchase := models.Purchase{
		GrBatchID:       body.GroupBatchesID,
		StatusPaymentID: 1, //pending
		PriceID:         price.ID,

		UserID:     &token.ID,
		URLConfirm: &confirmationCode,
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

	if err := db.
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

	if err := db.
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
