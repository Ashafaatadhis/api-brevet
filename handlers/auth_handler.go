package handlers

import (
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"new-brevet-be/validation"
	"time"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"

	log "github.com/sirupsen/logrus"
)

// Register untuk register user baru
func Register() fiber.Handler {
	db := config.DB
	return func(c *fiber.Ctx) error {
		body := c.Locals("body").(validation.UserRegister)
		newLog := log.WithFields(log.Fields{

			"event": "register",
		})
		// Hash password
		hashedPassword, err := utils.HashPassword(body.Password)
		if err != nil {

			newLog.Error("ERROR: Failed to hash password:", err.Error())
			return utils.Response(c, fiber.StatusInternalServerError, "Invalid server errror", nil, nil, nil)

		}

		body.RoleID = 4
		// Simpan user ke database
		body.Password = hashedPassword
		if err := db.Create(&body).Error; err != nil {
			newLog.Warn("WARNING: Failed to register user")
			return utils.Response(c, fiber.StatusBadRequest, "Failed to register user", nil, nil, nil)

		}

		// Mengambil user dengan preload role untuk mendapatkan data lengkap
		var userWithRole dto.ResponseUser
		if err := db.Preload("Role").First(&userWithRole, body.ID).Error; err != nil {

			newLog.Error("ERROR: Failed to fetch user with role:", err.Error())
			return utils.Response(c, fiber.StatusInternalServerError, "Failed to register user", nil, nil, nil)
		}

		newLog.Info("User registered successfully")
		return utils.Response(c, fiber.StatusOK, "User registered successfully", userWithRole, nil, nil)

	}
}

// Me adalah handle untuk mendapatkan data akun sesuai token
func Me(c *fiber.Ctx) error {
	db := config.DB
	user := c.Locals("user").(middlewares.User)
	newLog := log.WithFields(log.Fields{
		"user_id": user.ID,
		"event":   "get_profile",
	})
	var myUser models.User
	var userWithRole dto.ResponseUser

	// Mengambil data user dari database dengan preload untuk Role
	if err := db.Preload("Role").Preload("Profile").Preload("Profile.Golongan").First(&myUser, user.ID).Error; err != nil {

		newLog.Warn("WARNING: Failed to get info user")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get info user", nil, nil, nil)
	}

	// Menggunakan myUser yang sudah diisi dari DB untuk dipetakan ke userWithRole
	if err := dto_mapper.Map(&userWithRole, myUser); err != nil {
		newLog.Error("ERROR: Error during mapping:", err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map registration response", nil, nil, nil)
	}

	newLog.Info("Get Info successfully")
	return utils.Response(c, fiber.StatusOK, "Get User successfully", userWithRole, nil, nil)
}

// Login untuk login user dan menghasilkan JWT token
func Login() fiber.Handler {
	db := config.DB
	return func(c *fiber.Ctx) error {

		var user models.User
		newLog := log.WithFields(log.Fields{
			"event": "login",
		})
		body := c.Locals("body").(validation.UserLogin)

		// Cari user berdasarkan username dan preload role
		if err := db.Preload("Role").Where("username = ?", body.Username).First(&user).Error; err != nil {
			newLog.Warn("WARNING: User doesnt exist")
			return utils.Response(c, fiber.StatusUnauthorized, "Invalid username or password", nil, nil, nil)
		}

		// Verifikasi password
		if !utils.CheckPasswordHash(body.Password, user.Password) {
			newLog.Warn("WARNING: Wrong Password")
			return utils.Response(c, fiber.StatusUnauthorized, "Invalid username or password", nil, nil, nil)
		}

		// Generate JWT token
		token, err := utils.GenerateToken(user.ID, user.Username, user.Role.Name, user.RoleID)
		if err != nil {
			newLog.Error("ERROR: Failed to generate token:", err.Error())
			return utils.Response(c, fiber.StatusInternalServerError, "Invalid server error", nil, nil, nil)
		}

		var userWithRole dto.ResponseUser

		if err := dto_mapper.Map(&userWithRole, user); err != nil {

			newLog.Error("ERROR: Error during mapping:", err.Error())
			return utils.Response(c, fiber.StatusInternalServerError, "Failed to map user response", nil, nil, nil)
		}

		newLog.Info("Login Successfully")
		// Kirim token sebagai response
		return utils.Response(c, fiber.StatusOK, "Login successful", nil, userWithRole, &token)

	}
}

// Logout untuk menghapus token
func Logout() fiber.Handler {
	db := config.DB
	return func(c *fiber.Ctx) error {

		user := c.Locals("user").(middlewares.User)
		newLog := log.WithFields(log.Fields{
			"user_id": user.ID,
			"event":   "logout",
		})
		// Simpan token ke database dengan expiry
		expiration := time.Now().Add(time.Hour * 24) // Sesuaikan dengan masa berlaku token JWT
		blacklistToken := models.TokenBlacklist{
			Token:     user.Token,
			ExpiredAt: expiration,
		}

		if err := db.Create(&blacklistToken).Error; err != nil {
			newLog.Error("ERROR: failed to save token:", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save token",
			})
		}

		newLog.Info("Logout successfull")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logout successful",
		})
	}

}
