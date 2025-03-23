package handlers

import (
	"fmt"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

// GetMateriByID handler untuk method GET
func GetMateriByID(c *fiber.Ctx) error {
	log := logrus.WithFields(logrus.Fields{"event": "get_materi_by_id"})

	db := config.DB
	materiID := c.Params("materiId")

	var materi models.Materi

	// Cari pertemuan berdasarkan ID dan preload relasi
	if err := db.Where("id = ?", materiID).
		First(&materi).Error; err != nil {
		log.WithFields(logrus.Fields{"id": materiID}).WithError(err).Error("Failed to fetch materi by ID")
		return utils.Response(c, fiber.StatusNotFound, "materi not found", nil, nil, nil)
	}

	// Inisialisasi response
	var materiResponse dto.MateriResponse
	if err := copier.CopyWithOption(&materiResponse, &materi, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map materi response")
		return err
	}

	log.WithFields(logrus.Fields{"id": materiID}).Info("Successfully fetched materi by ID")
	return utils.Response(c, fiber.StatusOK, "materi fetched successfully", materiResponse, nil, nil)
}

// CreateMateri handler untuk POST materi
func CreateMateri(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "create_materi"})

	db := config.DB
	body := c.Locals("body").(dto.CreateMateriRequest)

	// Ambil pertemuanId dari parameter URL
	pertemuanIDStr := c.Params("pertemuanId")
	pertemuanID, err := strconv.Atoi(pertemuanIDStr)
	if err != nil {
		log.WithError(err).Error("Invalid pertemuan_id format")
		return utils.Response(c, fiber.StatusBadRequest, "Invalid pertemuan ID", nil, nil, nil)
	}

	var materi models.Materi

	if err := copier.CopyWithOption(&materi, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}

	fileURL, err := utils.UploadFile(c, "fileURL", "materi")

	if err != nil {
		log.WithError(err).Error("Failed Upload File")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed Upload File", nil, nil, nil)
	} else if fileURL != nil {
		materi.FileURL = *fileURL
	}

	materi.PertemuanID = pertemuanID

	// Jika file berhasil di-upload, simpan path-nya
	if fileURL != nil {
		materi.FileURL = *fileURL
	}

	if err := db.Create(&materi).Scan(&materi).Error; err != nil {
		log.WithError(err).Error("Failed to create materi")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to create materi", nil, nil, nil)
	}

	var materiResponse dto.MateriResponse

	if err := copier.CopyWithOption(&materiResponse, &materi, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{"materi_id": materi.ID}).Info("Materi created successfully")
	return utils.Response(c, fiber.StatusOK, "Materi created successfully", materiResponse, nil, nil)
}

// UpdateMateri handler untuk PUT materi
func UpdateMateri(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "update_materi"})

	db := config.DB
	body := c.Locals("body").(dto.UpdateMateriRequest)

	// Ambil materiId dari parameter URL
	materiIDStr := c.Params("materiId")
	materiID, err := strconv.Atoi(materiIDStr)
	if err != nil {
		log.WithError(err).Error("Invalid materi_id format")
		return utils.Response(c, fiber.StatusBadRequest, "Invalid materi ID", nil, nil, nil)
	}

	var materi models.Materi
	if err := db.First(&materi, materiID).Error; err != nil {
		log.WithError(err).Error("Materi not found")
		return utils.Response(c, fiber.StatusNotFound, "Materi not found", nil, nil, nil)
	}

	if err := copier.CopyWithOption(&materi, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}

	// Simpan path file lama untuk nanti dihapus jika ada file baru
	oldFileURL := materi.FileURL

	file, err := c.FormFile("fileURL")

	if file != nil { // Jika ada file yang dikirim dalam request
		fileURL, err := utils.UploadFile(c, "fileURL", "materi")
		if err != nil {
			log.WithError(err).Error("Failed to upload file")
			return utils.Response(c, fiber.StatusInternalServerError, "Failed to upload file", nil, nil, nil)
		} else if fileURL != nil {

			log.Print(*fileURL, "TEST filelURL")
			materi.FileURL = *fileURL // Update file jika ada
			log.Print(materi.FileURL, "TEST filelURL 2")
		}
	}

	// Update hanya field yang dikirim dalam request
	if err := db.Model(&materi).Updates(materi).Error; err != nil {
		log.WithError(err).Error("Failed to update materi")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update materi", nil, nil, nil)
	}

	// Hapus gambar lama jika ada
	if oldFileURL != "" {
		oldAvatarPath := fmt.Sprintf("./public/uploads/%s", oldFileURL) // Sesuaikan path
		if err := os.Remove(oldAvatarPath); err != nil {
			log.Warnf("Failed to delete old avatar: %s", err.Error())
		}
	}

	// Mapping ke response DTO
	var materiResponse dto.MateriResponse
	if err := copier.CopyWithOption(&materiResponse, &materi, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{"materi_id": materi.ID}).Info("Materi updated successfully")
	return utils.Response(c, fiber.StatusOK, "Materi updated successfully", materiResponse, nil, nil)
}

// DeleteMateri handler untuk Delete materi
func DeleteMateri(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "delete_materi"})

	db := config.DB

	// Ambil ID Materi dari parameter URL
	materiIDStr := c.Params("materiId")
	materiID, err := strconv.Atoi(materiIDStr)
	if err != nil {
		log.WithError(err).Error("Invalid materi_id format")
		return utils.Response(c, fiber.StatusBadRequest, "Invalid materi ID", nil, nil, nil)
	}

	// Cek apakah materi ada
	var materi models.Materi
	if err := db.First(&materi, materiID).Error; err != nil {
		log.WithError(err).Error("Materi not found")
		return utils.Response(c, fiber.StatusNotFound, "Materi not found", nil, nil, nil)
	}

	// **Hapus file yang terkait (jika ada)**
	if materi.FileURL != "" {
		oldAvatarPath := fmt.Sprintf("./public/uploads/%s", materi.FileURL) // Sesuaikan path
		if err := os.Remove(oldAvatarPath); err != nil {
			log.Warnf("Failed to delete old avatar: %s", err.Error())
		}
	}

	// **Hapus materi dari database**
	if err := db.Delete(&materi).Error; err != nil {
		log.WithError(err).Error("Failed to delete materi")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to delete materi", nil, nil, nil)
	}

	log.WithFields(logrus.Fields{"materi_id": materi.ID}).Info("Materi deleted successfully")
	return utils.Response(c, fiber.StatusOK, "Materi deleted successfully", nil, nil, nil)
}
