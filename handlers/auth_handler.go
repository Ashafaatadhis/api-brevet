package handlers

import (
	"log"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"new-brevet-be/validation"
	"time"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
)

// Register untuk register user baru
func Register() fiber.Handler {
	db := config.DB
	return func(c *fiber.Ctx) error {
		body := c.Locals("body").(validation.UserRegister)

		// Hash password
		hashedPassword, err := utils.HashPassword(body.Password)
		if err != nil {
			log.Println("Failed to hash password:", err)
			return utils.Response(c, fiber.StatusInternalServerError, "Invalid server errror", nil, nil, nil)

		}

		body.RoleID = 4
		// Simpan user ke database
		body.Password = hashedPassword
		if err := db.Create(&body).Error; err != nil {
			return utils.Response(c, fiber.StatusBadRequest, "Failed to register user", nil, nil, nil)

		}

		// Mengambil user dengan preload role untuk mendapatkan data lengkap
		var userWithRole dto.ResponseUser
		if err := db.Preload("Role").First(&userWithRole, body.ID).Error; err != nil {
			log.Println("Failed to fetch user with role:", err)
			return utils.Response(c, fiber.StatusInternalServerError, "Failed to register user", nil, nil, nil)
		}

		return utils.Response(c, fiber.StatusOK, "User registered successfully", userWithRole, nil, nil)

	}
}

// Me adalah handle untuk mendapatkan data akun sesuai token
func Me(c *fiber.Ctx) error {
	db := config.DB
	user := c.Locals("user").(middlewares.User)

	var myUser models.User
	var userWithRole dto.ResponseUser

	// Mengambil data user dari database dengan preload untuk Role
	if err := db.Preload("Role").Preload("Profile").Preload("Profile.Golongan").First(&myUser, user.ID).Error; err != nil {
		log.Println("Failed to fetch user with role:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get info user", nil, nil, nil)
	}

	// Menggunakan myUser yang sudah diisi dari DB untuk dipetakan ke userWithRole
	if err := dto_mapper.Map(&userWithRole, myUser); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map registration response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Get User successfully", userWithRole, nil, nil)
}

// Login untuk login user dan menghasilkan JWT token
func Login() fiber.Handler {
	db := config.DB
	return func(c *fiber.Ctx) error {
		var user models.User
		body := c.Locals("body").(validation.UserLogin)
		// Cari user berdasarkan username dan preload role
		if err := db.Preload("Role").Where("username = ?", body.Username).First(&user).Error; err != nil {
			return utils.Response(c, fiber.StatusUnauthorized, "Invalid username or password", nil, nil, nil)
		}

		// Verifikasi password
		if !utils.CheckPasswordHash(body.Password, user.Password) {
			return utils.Response(c, fiber.StatusUnauthorized, "Invalid username or password", nil, nil, nil)
		}

		// Generate JWT token
		token, err := utils.GenerateToken(user.ID, user.Username, user.Role.Name, user.RoleID)
		if err != nil {
			log.Println("Failed to generate token:", err)
			return utils.Response(c, fiber.StatusInternalServerError, "Invalid server error", nil, nil, nil)
		}

		var userWithRole dto.ResponseUser

		if err := dto_mapper.Map(&userWithRole, user); err != nil {
			log.Println("Error during mapping:", err)
			return utils.Response(c, fiber.StatusInternalServerError, "Failed to map user response", nil, nil, nil)
		}

		// Kirim token sebagai response
		return utils.Response(c, fiber.StatusOK, "Login successful", nil, userWithRole, &token)

	}
}

// Logout untuk menghapus token
func Logout() fiber.Handler {
	db := config.DB
	return func(c *fiber.Ctx) error {

		user := c.Locals("user").(middlewares.User)

		// Simpan token ke database dengan expiry
		expiration := time.Now().Add(time.Hour * 24) // Sesuaikan dengan masa berlaku token JWT
		blacklistToken := models.TokenBlacklist{
			Token:     user.Token,
			ExpiredAt: expiration,
		}

		if err := db.Create(&blacklistToken).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save token",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logout successful",
		})
	}

}
