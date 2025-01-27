package handlers

import (
	"log"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"new-brevet-be/validation"
	"strings"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
)

// GetAllBatchMappping adalah handler untuk route get batch
func GetAllBatchMappping(c *fiber.Ctx) error {
	db := config.DB

	// Ambil query parameters
	search := c.Query("q", "")            // Pencarian (default kosong)
	filter := c.Query("filter", "")       // Filter (e.g., Kursus, Teacher)
	selectFields := c.Query("select", "") // Field yang diinginkan (e.g., name, id)
	limit := c.QueryInt("limit", 10)      // Batas jumlah data (default 10)
	page := c.QueryInt("page", 1)         // Halaman (default 1)

	// Pagination offset
	offset := (page - 1) * limit

	var groupBatch []models.GroupBatch

	// Query builder
	query := db.Model(&models.GroupBatch{}).
		Preload("Kursus").
		Preload("Teacher").
		Preload("Batch")

	// Apply search query
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Apply filter
	if filter != "" {
		query = query.Where("kursus = ?", filter)
	}

	// Apply select fields
	if selectFields != "" {
		// Pisahkan field berdasarkan koma (e.g., "name,id")
		fields := strings.Split(selectFields, ",")
		query = query.Select(fields)
	}

	// Hitung total data sebelum pagination
	var totalData int64
	if err := query.Count(&totalData).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, false, "Failed to count total data", nil, nil, err.Error())
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&groupBatch).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, false, "Failed to get mapping batch", nil, nil, err.Error())
	}

	// Mapping ke DTO
	var batchResponseList []dto.GroupBatchResponse
	if err := dto_mapper.Map(&batchResponseList, groupBatch); err != nil {
		log.Println("Error during mapping:", err)
		return utils.NewResponse(c, fiber.StatusInternalServerError, false, "Failed to map batch response", nil, nil, err.Error())
	}

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	// Success response
	return utils.NewResponse(c, fiber.StatusOK, true, "Mapping batch retrieved successfully", batchResponseList, meta, nil)
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
