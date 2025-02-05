package handlers

import (
	"log"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/models"
	"new-brevet-be/utils"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// CreatePertemuan handler untuk POST pertemuan
func CreatePertemuan(c *fiber.Ctx) error {

	db := config.DB
	body := c.Locals("body").(dto.CreatePertemuanRequest)

	var pertemuan models.Pertemuan

	// Automapping
	if err := dto_mapper.Map(&pertemuan, body); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map pertemuan response", nil, nil, nil)
	}

	if err := db.Create(&pertemuan).Error; err != nil {
		log.Println("Error create:", err)
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create Pertemuan", nil, nil, nil)
	}

	if err := db.Where("id = ?", pertemuan.ID).
		Preload("GroupBatch").
		Preload("GroupBatch.Kursus").
		Preload("GroupBatch.Teacher").
		Preload("GroupBatch.Batch").
		First(&pertemuan).Error; err != nil {
		log.Println("Failed to fetch pertemuan with relations:", err)
		return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
	}

	// Inisialisasi response
	var pertemuanResponseList dto.PertemuanResponse

	// Automapping
	if err := dto_mapper.Map(&pertemuanResponseList, pertemuan); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map pertemuan response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Pertemuan created successfully", pertemuanResponseList, nil, nil)
}

// EditPertemuan handler untuk PUT pertemuan
func EditPertemuan(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(dto.EditPertemuanRequest)
	ID := c.Params("id")

	var pertemuan models.Pertemuan

	// 1️⃣ Pastikan data lama diambil sebelum update
	if err := db.Preload("GroupBatch").
		Preload("GroupBatch.Kursus").
		Preload("GroupBatch.Teacher").
		Preload("GroupBatch.Batch").First(&pertemuan, ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
		}
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
	}

	// 2️⃣ Cek jika name dan gr_batch_id tidak berubah
	if pertemuan.Name == utils.GetStringValue(body.Name) && pertemuan.GrBatchID == utils.GetIntValue(body.GrBatchID) {
		var pertemuanResponseList dto.PertemuanResponse

		if err := dto_mapper.Map(&pertemuanResponseList, pertemuan); err != nil {
			log.Println("Error during mapping:", err)
			return utils.Response(c, fiber.StatusInternalServerError, "Failed to map pertemuan response", nil, nil, nil)
		}

		return utils.Response(c, fiber.StatusOK, "Pertemuan updated successfully", pertemuanResponseList, nil, nil)

	}

	// 4️⃣ Lakukan update hanya jika tidak ada error
	if err := copier.CopyWithOption(&pertemuan, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}

	if err := db.Model(&pertemuan).Updates(pertemuan).Error; err != nil {
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update pertemuan", nil, nil, nil)
	}

	// 5️⃣ Ambil data terbaru setelah update
	if err := db.Where("id = ?", ID).
		Preload("GroupBatch").
		Preload("GroupBatch.Kursus").
		Preload("GroupBatch.Teacher").
		Preload("GroupBatch.Batch").
		First(&pertemuan).Error; err != nil {
		log.Println("Failed to fetch pertemuan with relations:", err)
		return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
	}

	// 6️⃣ Automapping response
	var pertemuanResponseList dto.PertemuanResponse
	if err := dto_mapper.Map(&pertemuanResponseList, pertemuan); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map pertemuan response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Pertemuan updated successfully", pertemuanResponseList, nil, nil)
}
