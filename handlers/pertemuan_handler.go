package handlers

import (
	"fmt"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/services"
	"new-brevet-be/utils"
	"strconv"
	"strings"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GetAllPertemuanByClass handler untuk method GET
func GetAllPertemuanByClass(c *fiber.Ctx) error {
	db := config.DB
	log := logrus.WithFields(logrus.Fields{"event": "get_all_pertemuan_by_class"})

	// Ambil parameter ID dari URL
	groupBatchID := c.Params("id")

	// Ambil query parameters
	search := c.Query("q", "")
	sort := c.Query("sort", "id")
	order := c.Query("order", "asc")
	selectFields := c.Query("select", "")
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)

	// Pagination offset
	offset := (page - 1) * limit

	log.WithField("group_batch_id", groupBatchID).Info("Fetching pertemuan for group batch")

	// Ambil valid sort fields dari tabel Pertemuan
	validSortFields, err := utils.GetValidSortFields(&models.Pertemuan{})
	if err != nil {
		log.WithError(err).Error("Failed to get valid sort fields")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get valid sort fields", nil, nil, nil)
	}

	// Validasi sorting field
	if !validSortFields[sort] {
		sort = "id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	var pertemuanList []models.Pertemuan
	query := db.Model(&models.Pertemuan{}).
		Where("gr_batch_id = ?", groupBatchID). // Filter berdasarkan GroupBatch ID
		Preload("Materis")

	// Apply search query
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Apply select fields
	if selectFields != "" {
		fields := strings.Split(selectFields, ",")
		query = query.Select(fields)
	}

	// Apply sorting
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	// Hitung total data sebelum pagination
	var totalData int64
	if err := query.Count(&totalData).Error; err != nil {
		log.WithError(err).Error("Failed to count total pertemuan")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to count total pertemuan", nil, nil, nil)
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&pertemuanList).Error; err != nil {
		log.WithError(err).Error("Failed to fetch pertemuan")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to fetch pertemuan", nil, nil, nil)
	}

	log.WithField("total_data", totalData).Info("Successfully fetched pertemuan for group batch")

	// Mapping ke DTO response
	var responseList []dto.PertemuanResponse
	if err := copier.CopyWithOption(&responseList, &pertemuanList, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map pertemuan response")
		return err
	}

	// Set struct kosong menjadi nil
	utils.TransformResponse(&responseList)

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "Pertemuan fetched successfully", responseList, meta, nil)
}

// GetPertemuanByClassByID handler untuk method GET
func GetPertemuanByClassByID(c *fiber.Ctx) error {
	log := logrus.WithFields(logrus.Fields{"event": "get_pertemuan_by_id"})

	db := config.DB
	pertemuanID := c.Params("pertemuanId")
	groupBatchID := c.Params("id")

	var pertemuan models.Pertemuan

	// Cari pertemuan berdasarkan ID dan preload relasi
	if err := db.Where("id = ? AND gr_batch_id = ?", pertemuanID, groupBatchID).
		Preload("Materis").
		First(&pertemuan).Error; err != nil {
		log.WithFields(logrus.Fields{"id": pertemuanID}).WithError(err).Error("Failed to fetch pertemuan by ID")
		return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
	}

	// Inisialisasi response
	var pertemuanResponse dto.PertemuanResponse
	if err := copier.CopyWithOption(&pertemuanResponse, &pertemuan, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map pertemuan response")
		return err
	}

	// Set struct kosong menjadi nil
	utils.TransformResponse(&pertemuanResponse)

	log.WithFields(logrus.Fields{"id": pertemuanID}).Info("Successfully fetched pertemuan by ID")
	return utils.Response(c, fiber.StatusOK, "Pertemuan fetched successfully", pertemuanResponse, nil, nil)
}

// CreatePertemuan handler untuk POST pertemuan
func CreatePertemuan(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "create_pertemuan"})

	db := config.DB
	body := c.Locals("body").(dto.CreatePertemuanRequest)

	// Ambil grBatchID dari parameter URL
	grBatchIDStr := c.Params("id")
	grBatchID, err := strconv.Atoi(grBatchIDStr)
	if err != nil {
		log.WithError(err).Error("Invalid group_batch_id format")
		return utils.Response(c, fiber.StatusBadRequest, "Invalid group_batch_id", nil, nil, nil)
	}

	// Cek apakah kombinasi unik
	isUnique, err := services.IsPertemuanUnique(grBatchID, body.Name, nil)
	if err != nil {
		log.WithError(err).Error("Failed to check uniqueness of pertemuan")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to validate uniqueness", nil, nil, nil)
	}

	if !isUnique {
		return utils.Response(c, fiber.StatusBadRequest, "Pertemuan dengan nama ini sudah ada dalam GroupBatch ini", nil, nil, nil)
	}

	var pertemuan models.Pertemuan

	// Automapping dari body ke struct `Pertemuan`
	if err := dto_mapper.Map(&pertemuan, body); err != nil {
		log.WithError(err).Error("Error during mapping")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map pertemuan response", nil, nil, nil)
	}

	// **Set `GrBatchID` secara manual**
	pertemuan.GrBatchID = grBatchID

	// Insert ke database
	if err := db.Create(&pertemuan).Error; err != nil {
		log.WithError(err).Error("Error creating pertemuan")
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create Pertemuan", nil, nil, nil)
	}

	// Ambil kembali data pertemuan dengan relasi
	if err := db.Where("id = ?", pertemuan.ID).
		Preload("Materis").
		First(&pertemuan).Error; err != nil {
		log.WithError(err).Error("Failed to fetch pertemuan with relations")
		return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
	}

	// Automapping response
	var pertemuanResponse dto.PertemuanResponse
	if err := dto_mapper.Map(&pertemuanResponse, pertemuan); err != nil {
		log.WithError(err).Error("Error during mapping response")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map pertemuan response", nil, nil, nil)
	}

	log.WithFields(logrus.Fields{"id": pertemuan.ID}).Info("Pertemuan created successfully")
	return utils.Response(c, fiber.StatusOK, "Pertemuan created successfully", pertemuanResponse, nil, nil)
}

// EditPertemuan handler untuk PUT pertemuan
func EditPertemuan(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "edit_pertemuan"})

	db := config.DB
	body := c.Locals("body").(dto.EditPertemuanRequest)

	pertemuanID := c.Params("pertemuanId")

	var pertemuan models.Pertemuan

	// Ambil data sebelum update
	if err := db.Preload("Materis").
		First(&pertemuan, pertemuanID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.WithFields(logrus.Fields{"id": pertemuanID}).Warn("Pertemuan not found")
			return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
		}
		log.WithError(err).Error("Failed to fetch pertemuan")
		return utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
	}

	// // Cek jika tidak ada perubahan
	// if pertemuan.Name == utils.GetStringValue(body.Name) && pertemuan.GrBatchID == utils.GetIntValue(body.GrBatchID) {
	// 	var pertemuanResponse dto.PertemuanResponse
	// 	if err := dto_mapper.Map(&pertemuanResponse, pertemuan); err != nil {
	// 		log.WithError(err).Error("Error during mapping response")
	// 		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map pertemuan response", nil, nil, nil)
	// 	}
	// 	log.WithFields(logrus.Fields{"id": ID}).Info("No changes detected, returning existing data")
	// 	return utils.Response(c, fiber.StatusOK, "Pertemuan updated successfully", pertemuanResponse, nil, nil)
	// }

	// Lakukan update
	if err := copier.CopyWithOption(&pertemuan, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Error copying data")
		return err
	}

	if err := db.Model(&pertemuan).Updates(pertemuan).Error; err != nil {
		log.WithError(err).Error("Failed to update pertemuan")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update pertemuan", nil, nil, nil)
	}

	// Ambil kembali data setelah update
	if err := db.Where("id = ?", pertemuanID).
		Preload("Materis").
		First(&pertemuan).Error; err != nil {
		log.WithError(err).Error("Failed to fetch updated pertemuan")
		return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
	}

	// Automapping response
	var pertemuanResponse dto.PertemuanResponse
	if err := dto_mapper.Map(&pertemuanResponse, pertemuan); err != nil {
		log.WithError(err).Error("Error during mapping response")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map pertemuan response", nil, nil, nil)
	}

	log.WithFields(logrus.Fields{"id": pertemuanID}).Info("Pertemuan updated successfully")
	return utils.Response(c, fiber.StatusOK, "Pertemuan updated successfully", pertemuanResponse, nil, nil)
}

// DeletePertemuan handler untuk DELETE pertemuan
func DeletePertemuan(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "delete_pertemuan"})

	db := config.DB
	pertemuanID := c.Params("pertemuanId")

	var pertemuan models.Pertemuan

	// Cek apakah pertemuan ada
	if err := db.First(&pertemuan, pertemuanID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.WithFields(logrus.Fields{"id": pertemuanID}).Warn("Pertemuan not found")
			return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
		}
		log.WithError(err).Error("Failed to fetch pertemuan")
		return utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
	}

	// Hapus pertemuan dari database
	if err := db.Delete(&pertemuan).Error; err != nil {
		log.WithError(err).Error("Failed to delete pertemuan")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to delete pertemuan", nil, nil, nil)
	}

	log.WithFields(logrus.Fields{"id": pertemuanID}).Info("Pertemuan deleted successfully")
	return utils.Response(c, fiber.StatusOK, "Pertemuan deleted successfully", nil, nil, nil)
}
