package handlers

import (
	"fmt"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"new-brevet-be/validation"
	"os"
	"strings"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
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
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{
		"user_id": token.ID,
		"event":   "create_manage_user",
	})
	// Hash password
	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		log.Error("ERROR: Failed to hash password:", err.Error())
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
		log.Error("ERROR: Failed to create user: ", err.Error())
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create user", nil, nil, nil)

	}

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
		log.Error("ERROR: Failed to create user profile: ", err.Error())
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create User", nil, nil, nil)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {

		log.Error("ERROR: Failed to commit transaction: ", err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to create user", nil, nil, nil)
	}

	// Mengambil user dengan preload role untuk mendapatkan data lengkap
	var userWithRole dto.ResponseUser
	if err := db.Preload("Role").Preload("Profile").Preload("Profile.Golongan").First(&user, user.ID).Error; err != nil {
		log.Error("ERROR: Failed to fetch user with role:", err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to create user", nil, nil, nil)
	}

	// Automapping
	if err := dto_mapper.Map(&userWithRole, user); err != nil {
		log.Error("ERROR: Error during mapping:", err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map kursus response", nil, nil, nil)
	}

	log.Info("User created successfully")
	return utils.Response(c, fiber.StatusOK, "User created successfully", userWithRole, nil, nil)
}

// GetManageUser handler untuk mengambil semua user kecuali admin
func GetManageUser(c *fiber.Ctx) error {
	db := config.DB
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{
		"user_id": token.ID,
		"event":   "get_manage_user",
	})
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
		log.Error("Failed to count total data: ", err.Error())
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to count total data", nil, nil, err.Error())
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&usersWithRoles).Error; err != nil {
		log.Error("Failed to get user: ", err.Error())
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get user", nil, nil, err.Error())
	}

	var userWithRoleList []dto.ResponseUser

	// Automapping
	if err := dto_mapper.Map(&userWithRoleList, usersWithRoles); err != nil {
		log.Info("Error during mapping: ", err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to mapping user response", nil, nil, nil)
	}

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	log.Info("User get successfully")
	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "User get successfully", userWithRoleList, meta, nil)

}

// GetDetailManageUser handler untuk mengambil semua user berdasarkan id kecuali admin
func GetDetailManageUser(c *fiber.Ctx) error {
	db := config.DB
	userID := c.Params("id")
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{
		"user_id": token.ID,
		"event":   "get_detail_manage_user",
	})
	// Mengambil user dengan preload role untuk mendapatkan data lengkap
	var userWithRole models.User
	if err := db.Preload("Role").
		Preload("Profile").Preload("Profile.Golongan").Where("id = ?", userID).
		First(&userWithRole).Error; err != nil {
		log.Error("Failed to fetch user with role:", err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get user", nil, nil, nil)
	}
	var userWithRoleList dto.ResponseUser

	// Automapping
	if err := dto_mapper.Map(&userWithRoleList, userWithRole); err != nil {
		log.Error("Error during mapping:", err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to mapping user response", nil, nil, nil)
	}

	log.Info("Uset get successfully")
	return utils.Response(c, fiber.StatusOK, "User get successfully", userWithRoleList, nil, nil)
}

// UpdateManageUser adalah handler untuk update data pengguna beserta profilnya
func UpdateManageUser(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(validation.UpdateManageUser)

	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{
		"user_id": token.ID,
		"event":   "update_manage_user",
	})

	// Ambil ID dari parameter route
	userID := c.Params("id")
	if userID == "" {
		log.Warn("User ID is required")
		return utils.Response(c, fiber.StatusBadRequest, "User ID is required", nil, nil, nil)
	}

	// Cari pengguna berdasarkan ID dengan preload Role & Profile
	var user models.User
	if err := db.Preload("Role").Preload("Profile").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("user not found: ", err.Error())
			return utils.Response(c, fiber.StatusNotFound, "User not found", nil, nil, nil)
		}
		log.Error(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
	}

	if err := c.BodyParser(&body); err != nil {
		log.Warn("Invalida request body: ", err.Error())
		return utils.Response(c, fiber.StatusBadRequest, "Invalid request body", nil, nil, nil)
	}

	// Salin nilai dari body ke user.Profile hanya jika field tidak nil
	if err := copier.CopyWithOption(&user, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.Error("Failed copy struct with copier: ", err.Error())
		return err
	}

	if err := copier.CopyWithOption(&user.Profile, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.Error("Failed copy struct with copier: ", err.Error())
		return err
	}

	// return utils.Response(c, fiber.StatusOK, "For", nil, nil, nil)

	if err := db.Model(&user).Updates(user).Error; err != nil {
		log.Error("Failed to update role_id: ", err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update role_id", nil, nil, nil)
	}
	if err := db.Model(&user.Profile).Updates(user.Profile).Error; err != nil {

		log.Error("Failed to update role_id: ", err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update role_id", nil, nil, nil)
	}

	var userWithRole dto.ResponseUser
	if err := db.Preload("Role").Preload("Profile").Preload("Profile.Golongan").First(&user, userID).Error; err != nil {
		log.Error("Failed to fetch user with role:", err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to create user", nil, nil, nil)
	}

	// Automapping
	if err := dto_mapper.Map(&userWithRole, user); err != nil {
		log.Error("Error during mapping:", err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map kursus response", nil, nil, nil)
	}

	log.Info("User updated successfully")
	return utils.Response(c, fiber.StatusOK, "User updated successfully", userWithRole, nil, nil)
}

// DeleteManageUser fungsi untuk handling manage-user method delete
func DeleteManageUser(c *fiber.Ctx) error {
	db := config.DB

	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{
		"user_id": token.ID,
		"event":   "delete_manage_user",
	})

	// Ambil ID dari parameter route
	userID := c.Params("id")
	if userID == "" {
		log.Warn("User ID is required")
		return utils.Response(c, fiber.StatusBadRequest, "User ID is required", nil, nil, nil)
	}

	// Cari pengguna berdasarkan ID dengan preload Role
	var user deleteUser
	if err := db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("user not found: ", err.Error())
			return utils.Response(c, fiber.StatusNotFound, "User not found", nil, nil, nil)
		}
		log.Error(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
	}

	// Hapus pengguna berdasarkan ID
	if err := db.Delete(&user).Error; err != nil {
		log.Error(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to delete user", nil, nil, nil)
	}

	// Hapus gambar lama jika ada
	if user.Avatar != "" {
		oldAvatarPath := fmt.Sprintf("./public/uploads/%s", user.Avatar) // Sesuaikan path
		if err := os.Remove(oldAvatarPath); err != nil {
			log.Warn("Failed to delete old avatar: ", err.Error())

		}
	}

	log.Info("User successfully deleted")
	// Berikan respon sukses jika berhasil
	return utils.Response(c, fiber.StatusOK, "User successfully deleted", nil, nil, nil)

}
