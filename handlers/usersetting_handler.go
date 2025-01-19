package handlers

import (
	"fmt"
	"log"
	"new-brevet-be/config"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"new-brevet-be/validation"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UpdateUserProfile untuk memperbarui data profil pengguna
func UpdateUserProfile() fiber.Handler {
	db := config.DB
	return func(c *fiber.Ctx) error {
		var user models.User

		token := c.Locals("user").(middlewares.User)
		body := c.Locals("body").(validation.UserSetting)

		// Cari pengguna berdasarkan ID dan preload relasi Role
		if err := db.Preload("Role").First(&user, token.ID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.Response(c, fiber.StatusNotFound, "User not found", nil, nil, nil)
			}
			log.Print(err.Error())
			return utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
		}
		// Periksa apakah ada file avatar yang diupload
		avatar, err := c.FormFile("avatar")

		if err != nil {
			log.Print(err.Error())
		}

		var data *string
		path := "profile"
		if avatar != nil {
			data, err = utils.UploadFileHandler(c, avatar, &path)
			if err != nil {
				return err
			}
		}

		// Hapus gambar lama jika ada
		if user.Avatar != "" {
			oldAvatarPath := fmt.Sprintf("./public/uploads/%s", user.Avatar) // Sesuaikan path
			if err := os.Remove(oldAvatarPath); err != nil {
				log.Printf("Failed to delete old avatar: %s", err.Error())
			}
		}

		if data != nil {
			user.Avatar = *data
		}

		// Perbarui data pengguna
		user.Name = body.Name
		user.Username = body.Username
		user.Nohp = body.Nohp
		user.Email = body.Email

		// Simpan perubahan
		if err := db.Save(&user).Error; err != nil {
			log.Print(err.Error())
			return utils.Response(c, fiber.StatusInternalServerError, "Failed to update user", nil, nil, nil)
		}
		return utils.Response(c, fiber.StatusOK, "User profile updated successfully", user, nil, nil)
	}
}

// ChangePassword handling untuk merubah password user
func ChangePassword(c *fiber.Ctx) error {
	var user models.User
	db := config.DB

	token := c.Locals("user").(middlewares.User)
	body := c.Locals("body").(validation.ChangePassword)
	// Cari pengguna berdasarkan ID
	if err := db.First(&user, token.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Response(c, fiber.StatusNotFound, "User not found", nil, nil, nil)
		}
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
	}

	// Verifikasi password
	if !utils.CheckPasswordHash(body.OldPassword, user.Password) {
		return utils.Response(c, fiber.StatusUnauthorized, "Invalid Password", nil, nil, nil)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(body.NewPassword)
	if err != nil {
		log.Println("Failed to hash password:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Invalid server errror", nil, nil, nil)

	}

	user.Password = hashedPassword
	// Simpan perubahan
	if err := db.Save(&user).Error; err != nil {
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to change password", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "User successfully change password", nil, nil, nil)

}
