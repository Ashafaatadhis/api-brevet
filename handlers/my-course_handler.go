package handlers

import (
	"fmt"
	"log"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"strings"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
)

// GetMyCourse adalah handler untuk route my-course
func GetMyCourse(c *fiber.Ctx) error {
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

	var myCourses []models.Purchase
	query := db.Model(&models.Purchase{}).
		Where("user_id = ? AND status_payment_id = ?", user.ID, 2 /* 2 is LUNAS */).
		Preload("GroupBatches").
		Preload("GroupBatches.Kursus").
		Preload("GroupBatches.Teacher").
		Preload("GroupBatches.Batch").
		Preload("JenisKursus")

		// Apply search query
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
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
	if err := query.Offset(offset).Limit(limit).Find(&myCourses).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get my course", nil, nil, err.Error())
	}

	var myCoursesList []dto.MyCourseResponse
	if err := dto_mapper.Map(&myCoursesList, myCourses); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map my courses response", nil, nil, nil)
	}

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "my courses retrieved successfully", myCoursesList, meta, nil)
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
