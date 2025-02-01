package handlers

import (
	"fmt"
	"log"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"new-brevet-be/validation"
	"os"
	"strings"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type deleteUser struct {
	ID     uint   `json:"id"`
	Avatar string `json:"avatar"`
}

// TableName untuk representasi ke table db
func (deleteUser) TableName() string {
	return "users"
}

// PostManageUser adalah handler untuk route post manage-user
func PostManageUser(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(validation.PostManageUser)

	// Hash password
	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		log.Println("Failed to hash password:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Invalid server errror", nil, nil, nil)

	}

	// Mulai transaction
	tx := db.Begin()

	user := models.User{
		Name:     body.Name,
		Username: body.Username,
		Nohp:     body.Nohp,
		RoleID:   body.RoleID,
		Email:    body.Email,
		Password: hashedPassword,
	}

	if err := tx.Create(&user).Scan(&user).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create user", nil, nil, nil)

	}

	log.Print(user, " woiii anjing")

	profile := models.Profile{
		GolonganID: nil,
		UserID:     &user.ID,
		Institusi:  body.Institusi,
		Asal:       body.Asal,
		TglLahir:   body.TglLahir,
		Alamat:     body.Alamat,
	}

	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create User", nil, nil, nil)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Failed to commit transaction:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to create user", nil, nil, nil)
	}

	// Mengambil user dengan preload role untuk mendapatkan data lengkap
	var userWithRole dto.ResponseUser
	if err := db.Preload("Role").Preload("Profile").Preload("Profile.Golongan").First(&user, user.ID).Error; err != nil {
		log.Println("Failed to fetch user with role:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to create user", nil, nil, nil)
	}

	// Automapping
	if err := dto_mapper.Map(&userWithRole, user); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map kursus response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "User created successfully", userWithRole, nil, nil)
}

// GetManageUser handler untuk mengambil semua user kecuali admin
func GetManageUser(c *fiber.Ctx) error {
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
	validSortFields, err := utils.GetValidSortFields(&models.User{})
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

	// Mengambil user dengan preload role untuk mendapatkan data lengkap
	var usersWithRoles []models.User
	query := db.Model(&models.User{}).Preload("Role").Preload("Profile")

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
	if err := query.Offset(offset).Limit(limit).Find(&usersWithRoles).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get user", nil, nil, err.Error())
	}

	var userWithRoleList []dto.ResponseUser

	// Automapping
	if err := dto_mapper.Map(&userWithRoleList, usersWithRoles); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to mapping user response", nil, nil, nil)
	}

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}
	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "User get successfully", userWithRoleList, meta, nil)

}

// GetDetailManageUser handler untuk mengambil semua user berdasarkan id kecuali admin
func GetDetailManageUser(c *fiber.Ctx) error {
	db := config.DB
	userID := c.Params("id")

	// Mengambil user dengan preload role untuk mendapatkan data lengkap
	var userWithRole models.User
	if err := db.Preload("Role").
		Preload("Profile").Preload("Profile.Golongan").Where("id = ? AND role_id != ?", userID, 1).
		First(&userWithRole).Error; err != nil {
		log.Println("Failed to fetch user with role:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get user", nil, nil, nil)
	}
	var userWithRoleList dto.ResponseUser

	// Automapping
	if err := dto_mapper.Map(&userWithRoleList, userWithRole); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to mapping user response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "User get successfully", userWithRoleList, nil, nil)
}

// UpdateManageUser adalah handler untuk update data pengguna beserta profilnya
func UpdateManageUser(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(validation.UpdateManageUser)

	// Ambil ID dari parameter route
	userID := c.Params("id")
	if userID == "" {
		return utils.Response(c, fiber.StatusBadRequest, "User ID is required", nil, nil, nil)
	}

	// Cari pengguna berdasarkan ID dengan preload Role & Profile
	var user models.User
	if err := db.Preload("Role").Preload("Profile").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Response(c, fiber.StatusNotFound, "User not found", nil, nil, nil)
		}
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
	}

	if err := c.BodyParser(&body); err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "Invalid request body", nil, nil, nil)
	}

	// Salin nilai dari body ke user.Profile hanya jika field tidak nil
	if err := copier.CopyWithOption(&user, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}

	if err := copier.CopyWithOption(&user.Profile, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}

	return utils.Response(c, fiber.StatusOK, "For you jidan sayangkuh", user, nil, nil)

	if err := db.Model(&user).Updates(user).Error; err != nil {
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update role_id", nil, nil, nil)
	}
	if err := db.Model(&user.Profile).Updates(user.Profile).Error; err != nil {
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update role_id", nil, nil, nil)
	}

	var userWithRole dto.ResponseUser
	if err := db.Preload("Role").Preload("Profile").Preload("Profile.Golongan").First(&user, userID).Error; err != nil {
		log.Println("Failed to fetch user with role:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to create user", nil, nil, nil)
	}

	// Automapping
	if err := dto_mapper.Map(&userWithRole, user); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map kursus response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "User updated successfully", userWithRole, nil, nil)
}

// DeleteManageUser fungsi untuk handling manage-user method delete
func DeleteManageUser(c *fiber.Ctx) error {
	db := config.DB

	// Ambil ID dari parameter route
	userID := c.Params("id")
	if userID == "" {
		return utils.Response(c, fiber.StatusBadRequest, "User ID is required", nil, nil, nil)
	}

	// Cari pengguna berdasarkan ID dengan preload Role
	var user deleteUser
	if err := db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Response(c, fiber.StatusNotFound, "User not found", nil, nil, nil)
		}
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
	}

	// Hapus pengguna berdasarkan ID
	if err := db.Delete(&user).Error; err != nil {
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to delete user", nil, nil, nil)
	}

	// Hapus gambar lama jika ada
	if user.Avatar != "" {
		oldAvatarPath := fmt.Sprintf("./public/uploads/%s", user.Avatar) // Sesuaikan path
		if err := os.Remove(oldAvatarPath); err != nil {
			log.Printf("Failed to delete old avatar: %s", err.Error())
		}
	}

	// Berikan respon sukses jika berhasil
	return utils.Response(c, fiber.StatusOK, "User successfully deleted", nil, nil, nil)

}
