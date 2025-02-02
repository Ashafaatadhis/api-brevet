package handlers

import (
	"fmt"
	"log"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"strings"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
)

// GetAllRoles adalah handler untuk route roles
func GetAllRoles(c *fiber.Ctx) error {
	db := config.DB

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
	validSortFields, err := utils.GetValidSortFields(&models.Role{})
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

	var roles []models.Role
	query := db.Model(&models.Role{})
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
	if err := query.Offset(offset).Limit(limit).Find(&roles).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get roles", nil, nil, err.Error())
	}

	var rolesList []dto.RoleResponse
	if err := dto_mapper.Map(&rolesList, roles); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map roles response", nil, nil, nil)
	}

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "roles retrieved successfully", rolesList, meta, nil)
}

// GetDetailRole adalah handler untuk route roles/:id
func GetDetailRole(c *fiber.Ctx) error {
	db := config.DB

	ID := c.Params("id")

	var role models.Role
	if err := db.
		First(&role, ID).Error; err != nil {
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get role", nil, nil, nil)
	}

	var roleList dto.RoleResponse
	if err := dto_mapper.Map(&roleList, role); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map role response", nil, nil, nil)
	}
	return utils.Response(c, fiber.StatusOK, "role retrieved successfully", roleList, nil, nil)
}
