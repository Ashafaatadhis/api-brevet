package handlers

import (
	"log"
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

	// Mengambil semua batch dengan preload semua relasi
	var batchList []models.Batch
	if err := db.Preload("GroupBatches").
		Preload("GroupBatches.Kursus").
		Find(&batchList).Error; err != nil {
		log.Println("Failed to fetch batch with relations:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get batch", nil, nil, nil)
	}

	// Inisialisasi response
	var batchResponseList []dto.BatchResponse

	// Automapping
	if err := dto_mapper.Map(&batchResponseList, batchList); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	// Kembalikan response
	return utils.Response(c, fiber.StatusOK, "Batch retrieved successfully", batchResponseList, nil, nil)
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
