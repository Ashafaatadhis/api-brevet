package handlers

import (
	"fmt"
	"log"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"new-brevet-be/validation"

	dto_mapper "github.com/dranikpg/dto-mapper" // Impor dengan alias
	"github.com/gofiber/fiber/v2"
)

// GetBatch handler untuk mengambil semua batch dengan preload semua relasi
func GetBatch(c *fiber.Ctx) error {
	db := config.DB

	// Ambil query parameters
	search := c.Query("q", "")       // Pencarian berdasarkan nama
	sort := c.Query("sort", "id")    // Sorting field (default "id")
	order := c.Query("order", "asc") // Urutan sorting (default "asc")
	limit := c.QueryInt("limit", 10) // Batas jumlah data (default 10)
	page := c.QueryInt("page", 1)    // Halaman (default 1)

	// Pagination offset
	offset := (page - 1) * limit

	// Ambil valid sort fields secara otomatis dari tabel
	validSortFields, err := utils.GetValidSortFields(&models.Batch{})
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

	// Mengambil semua batch
	var batchList []models.Batch
	query := db.Model(&models.Batch{}).Preload("GroupBatches").Preload("GroupBatches.Kursus").Preload("GroupBatches.Teacher")

	// Apply search query
	if search != "" {
		query = query.Where("judul LIKE ?", "%"+search+"%")
	}

	// Apply sorting
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	// Hitung total data sebelum pagination
	var totalData int64
	if err := query.Count(&totalData).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to count total data", nil, nil, err.Error())
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&batchList).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get batch", nil, nil, err.Error())
	}

	// Inisialisasi response
	var batchResponseList []dto.BatchResponse

	// Automapping
	if err := dto_mapper.Map(&batchResponseList, batchList); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "Batch retrieved successfully", batchResponseList, meta, nil)
}

// GetDetailBatch handler untuk mengambil detail batch
func GetDetailBatch(c *fiber.Ctx) error {
	db := config.DB
	batchID := c.Params("id")

	// Mengambil kursus berdasarkan ID dengan preload semua relasi
	var batch models.Batch
	if err := db.Where("id = ?", batchID).
		Preload("GroupBatches").
		Preload("GroupBatches.Kursus").
		Preload("GroupBatches.Teacher").
		First(&batch).Error; err != nil {
		log.Println("Failed to fetch batch with relations:", err)
		return utils.Response(c, fiber.StatusNotFound, "Batch not found", nil, nil, nil)
	}

	// Inisialisasi response
	var batchResponseList dto.BatchResponse

	// Automapping
	if err := dto_mapper.Map(&batchResponseList, batch); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Batch retrieved successfully", batchResponseList, nil, nil)
}

// PostBatch adalah handler untuk route post batch
func PostBatch(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(validation.PostBatch)

	var batch models.Batch

	// Automapping
	if err := dto_mapper.Map(&batch, body); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	if err := db.Create(&batch).Error; err != nil {

		return utils.Response(c, fiber.StatusBadRequest, "Failed to create Batch", nil, nil, nil)
	}

	if err := db.Where("id = ?", batch.ID).
		Preload("GroupBatches").
		Preload("GroupBatches.Kursus").
		First(&batch).Error; err != nil {
		log.Println("Failed to fetch batch with relations:", err)
		return utils.Response(c, fiber.StatusNotFound, "Batch not found", nil, nil, nil)
	}

	// Inisialisasi response
	var batchResponseList dto.BatchResponse

	// Automapping
	if err := dto_mapper.Map(&batchResponseList, batch); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Batch created successfully", batchResponseList, nil, nil)
}

// UpdateBatch adalah handler untuk route update batch
func UpdateBatch(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(validation.PostBatch)
	batchID := c.Params("id")

	var batch models.Batch
	if err := db.First(&batch, batchID).Error; err != nil {
		return utils.Response(c, fiber.StatusNotFound, "Batch not found", nil, nil, nil)
	}

	// Automapping
	if err := dto_mapper.Map(&batch, body); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	if err := db.Save(&batch).Error; err != nil {

		return utils.Response(c, fiber.StatusBadRequest, "Failed to update Batch", nil, nil, nil)
	}

	if err := db.Where("id = ?", batch.ID).
		Preload("GroupBatches").
		Preload("GroupBatches.Kursus").
		First(&batch).Error; err != nil {
		log.Println("Failed to fetch batch with relations:", err)
		return utils.Response(c, fiber.StatusNotFound, "Batch not found", nil, nil, nil)
	}

	// Inisialisasi response
	var batchResponseList dto.BatchResponse

	// Automapping
	if err := dto_mapper.Map(&batchResponseList, batch); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Batch updated successfully", batchResponseList, nil, nil)
}

// DeleteBatch adalah handler untuk route delete batch
func DeleteBatch(c *fiber.Ctx) error {
	db := config.DB
	batchID := c.Params("id")
	tx := db.Begin()

	var batch models.Batch
	// Fetch batch berdasarkan batchID
	if err := db.First(&batch, batchID).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusNotFound, "Batch not found", nil, nil, nil)
	}

	// Menghapus asosiasi Kursus melalui tabel pivot GroupBatch
	// Pastikan tabel pivot sudah benar didefinisikan dan ada hubungan antara batch dan kursus
	if err := tx.Model(&models.GroupBatch{}).
		Where("batch_id = ?", batch.ID).
		Delete(&models.GroupBatch{}).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to remove Kursus associations", nil, nil, nil)
	}

	// Menghapus batch itu sendiri
	if err := tx.Delete(&batch).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to delete Batch", nil, nil, nil)
	}

	// Commit transaksi
	tx.Commit()

	return utils.Response(c, fiber.StatusOK, "Batch deleted successfully", nil, nil, nil)
}
