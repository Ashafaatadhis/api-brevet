package handlers

import (
	"log"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"new-brevet-be/validation"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
)

// GetAllBatchMappping adalah handler untuk route get batch
func GetAllBatchMappping(c *fiber.Ctx) error {
	db := config.DB

	var groupBatch []models.GroupBatch
	if err := db.Preload("Kursus").Preload("Teacher").Preload("Batch").Find(&groupBatch).Error; err != nil {
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get mapping batch", nil, nil, nil)
	}

	var batchResponseList []dto.GroupBatchResponse
	if err := dto_mapper.Map(&batchResponseList, groupBatch); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}
	return utils.Response(c, fiber.StatusOK, "mapping batch retrieved successfully", batchResponseList, nil, nil)
}

// GetDetailBatchMappping adalah handler untuk route get batch/:id
func GetDetailBatchMappping(c *fiber.Ctx) error {
	db := config.DB
	ID := c.Params("id")
	var groupBatch models.GroupBatch

	if err := db.Preload("Kursus").Preload("Teacher").Preload("Batch").First(&groupBatch, ID).Error; err != nil {
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get mapping batch", nil, nil, nil)
	}

	var batchResponseList dto.GroupBatchResponse
	if err := dto_mapper.Map(&batchResponseList, groupBatch); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "mapping batch retrieved successfully", batchResponseList, nil, nil)
}

// CreateBatchMapping adalah handler untuk route update batch-mapping
func CreateBatchMapping(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(validation.CreateBatchMapping)

	// Fetch the batch
	var batch models.Batch
	if err := db.First(&batch, body.BatchID).Error; err != nil {

		return utils.Response(c, fiber.StatusNotFound, "Batch not found", nil, nil, nil)
	}

	// Process the kursusList and update the pivot table
	var kursus models.Kursus
	if err := db.First(&kursus, body.KursusID).Error; err != nil {

		return utils.Response(c, fiber.StatusBadRequest, "Invalid 'kursus' ID", nil, nil, nil)
	}

	// Remove old GroupBatch entries for the batch to replace them
	if err := db.Where("batch_id = ? AND kursus_id = ?", batch.ID, kursus.ID).First(&models.GroupBatch{}).Error; err == nil {

		return utils.Response(c, fiber.StatusBadRequest, "This Batch already mapping in kursus", nil, nil, nil)
	}

	// Now, upsert or replace the associations in GroupBatch (pivot table)

	groupBatch := models.GroupBatch{
		BatchID:  &batch.ID, // Assuming BatchID is not a pointer, adjust if needed
		KursusID: &kursus.ID,
	}

	// log.Print(groupBatch, " memwk")
	// return utils.Response(c, fiber.StatusBadRequest, "Failed to associate Kursus", groupBatch, nil, nil)

	// Save the new GroupBatch (this will insert into the pivot table)
	if err := db.Create(&groupBatch).Error; err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusBadRequest, "Failed to associate Kursus", nil, nil, nil)
	}

	// Fetch the updated batch with associations
	if err := db.Where("id = ?", groupBatch.ID).Preload("Kursus").Preload("Teacher").Preload("Batch").First(&groupBatch).Error; err != nil {
		log.Println("Failed to fetch groupbatch with relations:", err)
		return utils.Response(c, fiber.StatusNotFound, "Batch not found", nil, nil, nil)
	}
	var batchResponseList dto.GroupBatchResponse
	if err := dto_mapper.Map(&batchResponseList, groupBatch); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}
	// Prepare response data

	return utils.Response(c, fiber.StatusOK, "Batch Mapping created successfully", batchResponseList, nil, nil)
}

// EditBatchMapping adalah handler untuk route update batch-mapping
func EditBatchMapping(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(validation.CreateBatchMapping)
	ID := c.Params("id")

	var groupBatch models.GroupBatch

	if err := db.First(&groupBatch, ID).Error; err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "This Mapping not exist", nil, nil, nil)
	}

	// Fetch the batch
	var batch models.Batch
	if err := db.First(&batch, body.BatchID).Error; err != nil {

		return utils.Response(c, fiber.StatusNotFound, "Batch not found", nil, nil, nil)
	}

	// Process the kursusList and update the pivot table
	var kursus models.Kursus
	if err := db.First(&kursus, body.KursusID).Error; err != nil {

		return utils.Response(c, fiber.StatusBadRequest, "Invalid 'kursus' ID", nil, nil, nil)
	}

	// Now, upsert or replace the associations in GroupBatch (pivot table)
	groupBatch.BatchID = &batch.ID
	groupBatch.KursusID = &kursus.ID

	// Save the new GroupBatch (this will insert into the pivot table)
	if err := db.Save(&groupBatch).Error; err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "Failed to associate Kursus", nil, nil, nil)
	}

	// Fetch the updated batch with associations
	if err := db.Where("id = ?", groupBatch.ID).Preload("Kursus").Preload("Teacher").Preload("Batch").First(&groupBatch).Error; err != nil {
		log.Println("Failed to fetch groupbatch with relations:", err)
		return utils.Response(c, fiber.StatusNotFound, "Batch not found", nil, nil, nil)
	}

	var batchResponseList dto.GroupBatchResponse
	if err := dto_mapper.Map(&batchResponseList, groupBatch); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	// Prepare response data

	return utils.Response(c, fiber.StatusOK, "Batch Mapping updated successfully", batchResponseList, nil, nil)
}

// DeleteBatchMapping adalah handler untuk route delete batch
func DeleteBatchMapping(c *fiber.Ctx) error {
	db := config.DB
	ID := c.Params("id")

	var groupBatch models.GroupBatch

	if err := db.First(&groupBatch, ID).Error; err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "This Mapping not exist", nil, nil, nil)
	}

	// Menghapus asosiasi Kursus melalui tabel pivot GroupBatch
	if err := db.Model(&groupBatch).
		Delete(&groupBatch, ID).Error; err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "Failed to remove Kursus associations", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Batch mapping deleted successfully", nil, nil, nil)
}
