package handlers

import (
	"log"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/models"
	"new-brevet-be/utils"

	dto_mapper "github.com/dranikpg/dto-mapper"
	"github.com/gofiber/fiber/v2"
)

type getKursus struct {
	HargaAsli   float64 `json:"harga_asli"`
	HargaDiskon float64 `json:"harga_diskon"`
}
type getGroupBatch struct {
	ID uint `json:"id"`
}

// GetAllRegistration adalah handler untuk route registration
func GetAllRegistration(c *fiber.Ctx) error {
	db := config.DB

	var users []models.User

	var registrationResponseList []dto.RegistrationResponse

	if err := db.Where("role_id = ?", 4).
		Preload("Profile").Find(&users).Error; err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "Failed to get all Registration", nil, nil, nil)
	}

	if err := dto_mapper.Map(&registrationResponseList, users); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map registration response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "success retrieved all registrations", registrationResponseList, nil, nil)
}

// GetDetailRegistration adalah handler untuk route registration
func GetDetailRegistration(c *fiber.Ctx) error {
	db := config.DB
	ID := c.Params("id")
	var users models.User

	var registrationResponseList dto.RegistrationResponse

	if err := db.Where("role_id = ? AND id = ?", 4, ID).
		Preload("Profile").Find(&users).Error; err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "Failed to get all Registration", nil, nil, nil)
	}

	if err := dto_mapper.Map(&registrationResponseList, users); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map registration response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "success retrieved all registrations", registrationResponseList, nil, nil)
}

// CreateRegistration adalah handler untuk route registration
func CreateRegistration(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(dto.CreateRegistrationRequest)

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
		RoleID:   4,
		Email:    body.Email,
		Password: hashedPassword,
	}

	// Simpan registration ke database
	if err := tx.Create(&user).Scan(&user).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create User", nil, nil, nil)
	}

	profile := models.Profile{
		GolonganID: nil,
		UserID:     &user.ID,
		Institusi:  body.Institusi,
		Asal:       body.Asal,
		TglLahir:   body.TglLahir,
		Alamat:     body.Alamat,
	}

	// upload gambar
	path := "bukti"

	dataNim, err := utils.UploadFile(c, "bukti_nim", path)
	if err != nil {
		log.Println("Failed to upload Bukti NIM:", err)
	} else if dataNim != nil {
		profile.BuktiNIM = dataNim
	}

	dataNik, err := utils.UploadFile(c, "bukti_nik", path)
	if err != nil {
		log.Println("Failed to upload Bukti NIK:", err)
	} else if dataNik != nil {
		profile.BuktiNIK = dataNik
	}

	if err := tx.Create(&profile).Scan(&profile).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create User", nil, nil, nil)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Failed to commit transaction:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to create Registration", nil, nil, nil)
	}

	// Automapping ke response
	var registrationResponseList dto.RegistrationResponse
	if err := dto_mapper.Map(&registrationResponseList, user); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map registration response", nil, nil, nil)
	}

	if err := dto_mapper.Map(&registrationResponseList.Profile, profile); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map registration response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Registration created successfully", registrationResponseList, nil, nil)
}

// EditRegistration adalah handler untuk route registration/:id (only status payment can edit)
func EditRegistration(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(dto.EditRegistrationRequest)
	userID := c.Params("id")

	var user models.User
	var profile models.Profile

	user.RoleID = 4
	// Mulai transaction
	tx := db.Begin()

	// Preload data user dengan profilnya
	if err := tx.First(&user, "id = ?", userID).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusNotFound, "User not found", nil, nil, nil)
	}

	// Ambil profil user berdasarkan user_id
	if err := tx.First(&profile, "user_id = ?", userID).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Profile not found", nil, nil, nil)
	}

	// Update GolonganID hanya jika body.GolonganID tidak nil
	if body.GolonganID != nil {
		if profile.GolonganID == nil {
			profile.GolonganID = new(int) // Alokasikan memori untuk pointer
		}
		*profile.GolonganID = *body.GolonganID // Menyalin pointer langsung
	}

	// Simpan perubahan profil ke database
	if err := tx.Save(&profile).Error; err != nil {
		tx.Rollback()
		log.Print(err.Error())
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update profile", nil, nil, nil)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Failed to commit transaction:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update registration", nil, nil, nil)
	}

	if err := db.Preload("Profile").First(&user, "id = ?", userID).Error; err != nil {

		return utils.Response(c, fiber.StatusNotFound, "User not found", nil, nil, nil)
	}

	// Automapping ke response
	var registrationResponseList dto.RegistrationResponse
	if err := dto_mapper.Map(&registrationResponseList, user); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map registration response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Registration updated successfully", registrationResponseList, nil, nil)
}
