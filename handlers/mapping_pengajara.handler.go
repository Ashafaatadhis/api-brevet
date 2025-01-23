package handlers

import (
	"log"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/models"
	"new-brevet-be/utils"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
)

// MappingPengajar adalah handler untuk route update batch-mapping
func MappingPengajar(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(dto.MappingPengajarRequest)
	ID := c.Params("id")

	var groupBatch models.GroupBatch

	if err := db.First(&groupBatch, ID).Error; err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "This Mapping not exist", nil, nil, nil)
	}

	// Now, upsert or replace the associations in GroupBatch (pivot table)

	groupBatch.TeacherID = &body.TeacherID

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

	return utils.Response(c, fiber.StatusOK, "Teacher Mapping updated successfully", batchResponseList, nil, nil)
}
