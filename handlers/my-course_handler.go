package handlers

import (
	"log"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
)

// GetMyCourse adalah handler untuk route my-course
func GetMyCourse(c *fiber.Ctx) error {
	db := config.DB
	user := c.Locals("user").(middlewares.User)

	var myCourses []models.Purchase
	if err := db.
		Where("user_id = ? AND status_payment_id = ?", user.ID, 2 /* 2 is LUNAS */).
		Preload("GroupBatches").
		Preload("GroupBatches.Kursus").
		Preload("GroupBatches.Teacher").
		Preload("GroupBatches.Batch").
		Preload("JenisKursus").
		Find(&myCourses).Error; err != nil {
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get my courses", nil, nil, nil)
	}

	var myCoursesList []dto.MyCourseResponse
	if err := dto_mapper.Map(&myCoursesList, myCourses); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map my courses response", nil, nil, nil)
	}
	return utils.Response(c, fiber.StatusOK, "my courses retrieved successfully", myCoursesList, nil, nil)
}

// GetMyCourseByID adalah handler untuk route my-course/:id
func GetMyCourseByID(c *fiber.Ctx) error {
	db := config.DB
	user := c.Locals("user").(middlewares.User)
	ID := c.Params("id")

	var myCourses models.Purchase
	if err := db.
		Where("user_id = ? AND status_payment_id = ?", user.ID, 2 /* 2 is LUNAS */).
		Preload("GroupBatches").
		Preload("GroupBatches.Kursus").
		Preload("GroupBatches.Teacher").
		Preload("GroupBatches.Batch").
		Preload("JenisKursus").
		First(&myCourses, ID).Error; err != nil {
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get my courses", nil, nil, nil)
	}

	var myCoursesList dto.MyCourseResponse
	if err := dto_mapper.Map(&myCoursesList, myCourses); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map my course response", nil, nil, nil)
	}
	return utils.Response(c, fiber.StatusOK, "my course retrieved successfully", myCoursesList, nil, nil)
}
