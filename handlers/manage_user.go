package handlers

import (
	"fmt"
	"log"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"new-brevet-be/validation"
	"os"

	"github.com/gofiber/fiber/v2"
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

	body.Password = hashedPassword
	if err := db.Create(&body).Error; err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create user", nil, nil, nil)

	}

	// Mengambil user dengan preload role untuk mendapatkan data lengkap
	var userWithRole dto.ResponseUser
	if err := db.Preload("Role").First(&userWithRole, body.ID).Error; err != nil {
		log.Println("Failed to fetch user with role:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to create user", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "User created successfully", userWithRole, nil, nil)
}

// GetManageUser handler untuk mengambil semua user kecuali admin
func GetManageUser(c *fiber.Ctx) error {
	db := config.DB
	// Mengambil user dengan preload role untuk mendapatkan data lengkap
	var usersWithRoles []models.User
	if err := db.Where("role_id != ?", 1).Preload("Role").Find(&usersWithRoles).Error; err != nil {
		log.Println("Failed to fetch users with role:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get users", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "User get successfully", usersWithRoles, nil, nil)
}

// GetDetailManageUser handler untuk mengambil semua user berdasarkan id kecuali admin
func GetDetailManageUser(c *fiber.Ctx) error {
	db := config.DB
	userID := c.Params("id")

	// Mengambil user dengan preload role untuk mendapatkan data lengkap
	var userWithRole dto.ResponseUser
	if err := db.Where("id = ? AND role_id != ?", userID, 1).
		Preload("Role").
		First(&userWithRole).Error; err != nil {
		log.Println("Failed to fetch user with role:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get user", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "User get successfully", userWithRole, nil, nil)
}

// UpdateManageUser adalah handler untuk route manage-user
func UpdateManageUser(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(validation.UpdateManageUser)
	// Ambil ID dari parameter route
	userID := c.Params("id")
	if userID == "" {
		return utils.Response(c, fiber.StatusBadRequest, "User ID is required", nil, nil, nil)
	}

	// Cari pengguna berdasarkan ID dengan preload Role
	var user models.User
	if err := db.Preload("Role").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Response(c, fiber.StatusNotFound, "User not found", nil, nil, nil)
		}
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
	}

	if err := c.BodyParser(&body); err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "Invalid request body", nil, nil, nil)
	}

	// Perbarui data pengguna
	user.Name = body.Name
	user.Username = body.Username
	user.Nohp = body.Nohp
	user.Email = body.Email

	if err := db.Save(&user).Error; err != nil {
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update user", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "User updated successfully", user, nil, nil)
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
