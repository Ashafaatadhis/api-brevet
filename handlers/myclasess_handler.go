package handlers

import (
	"fmt"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

// GetMyClasses mengambil semua kelas yang diajar oleh guru
func GetMyClasses(c *fiber.Ctx) error {
	db := config.DB
	user := c.Locals("user").(middlewares.User)

	log := logrus.WithFields(logrus.Fields{
		"event":   "get_all_pertemuan",
		"user_id": user.ID,
	})

	// Ambil query parameters
	search := c.Query("q", "")
	sort := c.Query("sort", "id")
	order := c.Query("order", "asc")
	selectFields := c.Query("select", "")
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)

	// Pagination offset
	offset := (page - 1) * limit

	log.Info("Fetching my classes")

	// Ambil valid sort fields dari tabel
	validSortFields, err := utils.GetValidSortFields(&models.GroupBatch{})
	if err != nil {
		log.WithError(err).Error("Failed to get valid sort fields")
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get valid sort fields", nil, nil, err.Error())
	}

	// Validasi sorting field
	if !validSortFields[sort] {
		sort = "id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	var myClasses []models.GroupBatch
	query := db.Model(&models.GroupBatch{}).
		Where("teacher_id = ?", user.ID).
		Preload("Kursus").
		Preload("Batch").
		Preload("Batch.Jenis").
		Preload("Batch.Kelas").
		Preload("Teacher")

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
		log.WithError(err).Error("Failed to count total data")
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to count total data", nil, nil, err.Error())
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&myClasses).Error; err != nil {
		log.WithError(err).Error("Failed to get my classes")
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get my classes", nil, nil, err.Error())
	}

	log.Info("My classes retrieved successfully")

	// Mapping ke DTO response
	var myClassesList []dto.MyClasessResponse
	if err := copier.CopyWithOption(&myClassesList, &myClasses, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map my classes response")
		return err
	}

	// Set struct kosong menjadi nil
	utils.TransformResponse(&myClassesList)

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "My classes retrieved successfully", myClassesList, meta, nil)
}

// GetMyClassByID mengambil detail kelas tertentu yang diajar oleh guru
func GetMyClassByID(c *fiber.Ctx) error {
	db := config.DB
	user := c.Locals("user").(middlewares.User)
	classID := c.Params("id")

	log := logrus.WithFields(logrus.Fields{
		"event":   "get_all_pertemuan",
		"user_id": user.ID,
	})

	log.Info("Fetching class details")

	var myClass models.GroupBatch
	if err := db.Where("id = ? AND teacher_id = ?", classID, user.ID).
		Preload("Kursus").
		Preload("Batch").
		Preload("Batch.Jenis").
		Preload("Batch.Kelas").
		Preload("Teacher").
		First(&myClass).Error; err != nil {
		log.WithError(err).Error("Class not found")
		return utils.NewResponse(c, fiber.StatusNotFound, "Class not found", nil, nil, err.Error())
	}

	log.Info("Class retrieved successfully")

	// Mapping ke DTO response
	var myClassResponse dto.MyClasessResponse
	if err := copier.CopyWithOption(&myClassResponse, &myClass, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map class response")
		return err
	}

	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "Class retrieved successfully", myClassResponse, nil, nil)
}
