package handlers

import (
	"fmt"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/models"
	"new-brevet-be/utils"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// GetAllAbsensi handler untuk mengambil semua absensi
func GetAllAbsensi(c *fiber.Ctx) error {
	db := config.DB

	log := logrus.WithFields(logrus.Fields{
		"event": "get_all_absensi",
	})

	// Ambil query parameters
	search := c.Query("q", "")       // Pencarian berdasarkan nama
	sort := c.Query("sort", "id")    // Sorting field (default "id")
	order := c.Query("order", "asc") // Urutan sorting (default "asc")
	limit := c.QueryInt("limit", 10) // Batas jumlah data (default 10)
	page := c.QueryInt("page", 1)    // Halaman (default 1)

	// Pagination offset
	offset := (page - 1) * limit

	// Ambil valid sort fields secara otomatis dari tabel
	validSortFields, err := utils.GetValidSortFields(&models.Absensi{})
	if err != nil {
		log.Info("Failed to get valid sort fields: ", err.Error())
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
	var absensiList []models.Absensi
	query := db.Model(&models.Absensi{})

	// Apply search query
	if search != "" {
		query = query.Where("judul LIKE ?", "%"+search+"%")
	}

	// Apply sorting
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	// Hitung total data sebelum pagination
	var totalData int64
	if err := query.Count(&totalData).Error; err != nil {
		log.Error("Failed to count total data: ", err.Error())
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to count total data", nil, nil, err.Error())
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&absensiList).Error; err != nil {
		log.Error("Failed to get absensi: ", err.Error())
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get batch", nil, nil, err.Error())
	}

	// Inisialisasi response
	var absensiResponseList []dto.AbsensiResponse

	// Automapping
	if err := dto_mapper.Map(&absensiResponseList, absensiList); err != nil {
		log.Error("Error during mapping: ", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map absensi response", nil, nil, nil)
	}

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	log.Info("Absensi retrieved successfully")
	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "Batch retrieved successfully", absensiResponseList, meta, nil)
}
