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
	"github.com/jinzhu/copier"
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

	var users []models.User

	var registrationResponseList []dto.RegistrationResponse

	query := db.Model(&models.User{}).Where("role_id = ?", 4).
		Preload("Profile")

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
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get registrations", nil, nil, err.Error())
	}

	if err := dto_mapper.Map(&registrationResponseList, users); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map registration response", nil, nil, nil)
	}

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "success retrieved all registrations", registrationResponseList, meta, nil)

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
	// Salin nilai dari body ke user.Profile hanya jika field tidak nil
	if err := copier.CopyWithOption(&user.Profile, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}
	// profile := models.Profile{
	// 	GolonganID: nil,
	// 	UserID:     &user.ID,
	// 	Institusi:  body.Institusi,
	// 	Asal:       body.Asal,
	// 	TglLahir:   body.TglLahir,
	// 	Alamat:     body.Alamat,
	// }

	// upload gambar
	path := "bukti"

	dataNim, err := utils.UploadFile(c, "bukti_nim", path)
	if err != nil {
		log.Println("Failed to upload Bukti NIM:", err)
	} else if dataNim != nil {
		user.Profile.BuktiNIM = dataNim
	}

	dataNik, err := utils.UploadFile(c, "bukti_nik", path)
	if err != nil {
		log.Println("Failed to upload Bukti NIK:", err)
	} else if dataNik != nil {
		user.Profile.BuktiNIK = dataNik
	}

	if err := tx.Create(&user.Profile).Scan(&user.Profile).Error; err != nil {
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

	if err := dto_mapper.Map(&registrationResponseList.Profile, user.Profile); err != nil {
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
