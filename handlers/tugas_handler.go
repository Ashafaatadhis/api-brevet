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

// GetTugasByID handler untuk method GET
func GetTugasByID(c *fiber.Ctx) error {
	log := logrus.WithFields(logrus.Fields{"event": "get_tugas_by_id"})
	token := c.Locals("user").(middlewares.User)
	db := config.DB
	tugasID := c.Params("tugasId")

	var tugas models.Tugas

	// Cari pertemuan berdasarkan ID dan preload relasi
	if err := db.Where("id = ?", tugasID).
		Preload("TugasFile").
		Preload("Jawaban", "user_id = ?", token.ID).
		First(&tugas).Error; err != nil {
		log.WithFields(logrus.Fields{"id": tugasID}).WithError(err).Error("Failed to fetch tugas by ID")
		return utils.Response(c, fiber.StatusNotFound, "Tugas not found", nil, nil, nil)
	}

	// Inisialisasi response
	var tugasResponse dto.TugasResponse
	if err := copier.CopyWithOption(&tugasResponse, &tugas, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map tugas response")
		return err
	}

	log.WithFields(logrus.Fields{"id": tugasID}).Info("Successfully fetched tugas by ID")
	return utils.Response(c, fiber.StatusOK, "Tugas fetched successfully", tugasResponse, nil, nil)
}

// CreateTugas handler untuk POST tugas
func CreateTugas(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "create_tugas"})

	db := config.DB
	body := c.Locals("body").(dto.CreateTugasRequest)

	// Ambil pertemuanId dari parameter URL
	pertemuanIDStr := c.Params("pertemuanId")
	pertemuanID, err := strconv.Atoi(pertemuanIDStr)
	if err != nil {
		log.WithError(err).Error("Invalid pertemuan_id format")
		return utils.Response(c, fiber.StatusBadRequest, "Invalid pertemuan ID", nil, nil, nil)
	}

	// Start Transaction
	tx := db.Begin()

	// Copy data dari request ke struct tugas
	var tugas models.Tugas
	if err := copier.CopyWithOption(&tugas, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		tx.Rollback()
		return err
	}
	tugas.PertemuanID = pertemuanID

	// Simpan tugas ke database
	if err := tx.Create(&tugas).Error; err != nil {
		log.WithError(err).Error("Failed to create tugas")
		tx.Rollback()
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to create tugas", nil, nil, nil)
	}

	// **Handle Multiple File Upload**
	fileURLs, err := utils.UploadMultipleFiles(c, "files", "tugas")
	if err != nil {
		log.WithError(err).Error("Failed to upload files")
		tx.Rollback()
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to upload files", nil, nil, nil)
	}

	// **Gunakan Association untuk Simpan File**
	var tugasFiles []models.TugasFile
	for _, fileURL := range fileURLs {
		tugasFiles = append(tugasFiles, models.TugasFile{
			FileURL: fileURL,
		})
	}

	if err := tx.Model(&tugas).Association("TugasFile").Append(&tugasFiles); err != nil {
		tx.Rollback()
		log.WithError(err).Error("Failed to associate files")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to associate files", nil, nil, nil)
	}

	// Commit transaction jika semua berhasil
	tx.Commit()

	// Ambil tugas yang baru dibuat dengan Preload tugas_files
	if err := db.Preload("TugasFile").First(&tugas, tugas.ID).Error; err != nil {
		log.WithError(err).Error("Failed to retrieve tugas with files")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to retrieve tugas", nil, nil, nil)
	}

	// Konversi ke Response DTO
	var tugasResponse dto.TugasResponse
	if err := copier.CopyWithOption(&tugasResponse, &tugas, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{"tugas_id": tugas.ID}).Info("Tugas created successfully")
	return utils.Response(c, fiber.StatusOK, "Tugas created successfully", tugasResponse, nil, nil)
}

// UpdateTugas handler untuk PATCH tugas
func UpdateTugas(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "update_tugas"})

	db := config.DB
	body := c.Locals("body").(dto.UpdateTugasRequest)

	// Ambil tugasId dari parameter URL
	tugasIDStr := c.Params("tugasId")
	tugasID, err := strconv.Atoi(tugasIDStr)
	if err != nil {
		log.WithError(err).Error("Invalid tugas_id format")
		return utils.Response(c, fiber.StatusBadRequest, "Invalid tugas ID", nil, nil, nil)
	}

	// Start Transaction
	tx := db.Begin()

	// **Ambil tugas lama beserta file-nya**
	var tugas models.Tugas
	if err := tx.Preload("TugasFile").First(&tugas, tugasID).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("Tugas not found")
		return utils.Response(c, fiber.StatusNotFound, "Tugas not found", nil, nil, nil)
	}

	// **Simpan daftar file lama sebelum dihapus**
	var oldFilePaths []string
	for _, oldFile := range tugas.TugasFile {
		oldFilePaths = append(oldFilePaths, fmt.Sprintf("./public/uploads/%s", oldFile.FileURL))
	}

	// **Update tugas dengan data baru**
	if err := copier.CopyWithOption(&tugas, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		tx.Rollback()
		return err
	}

	// **Gunakan `tx.Updates()` untuk update data tugas**
	if err := tx.Model(&tugas).Updates(tugas).Error; err != nil {
		log.WithError(err).Error("Failed to update tugas")
		tx.Rollback()
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update tugas", nil, nil, nil)
	}

	// **Handle file baru**
	fileURLs, err := utils.UploadMultipleFiles(c, "files", "tugas")
	if err != nil {
		log.WithError(err).Error("Failed to upload new files")
		tx.Rollback()
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to upload new files", nil, nil, nil)
	}

	// **Buat slice TugasFile baru**
	var tugasFiles []models.TugasFile
	for _, fileURL := range fileURLs {
		tugasFiles = append(tugasFiles, models.TugasFile{
			TugasID: &tugas.ID,
			FileURL: fileURL,
		})
	}

	// Hapus semua file lama dulu
	if err := tx.Where("tugas_id = ?", tugas.ID).Delete(&models.TugasFile{}).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("Failed to delete old files")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to delete old files", nil, nil, nil)
	}

	// Baru Replace dengan file yang baru
	if err := tx.Model(&tugas).Association("TugasFile").Replace(&tugasFiles); err != nil {
		tx.Rollback()
		log.WithError(err).Error("Failed to replace files")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to replace files", nil, nil, nil)
	}

	// **Commit transaction jika semua berhasil**
	tx.Commit()

	// **Hapus file lama dari storage setelah commit**
	for _, oldFilePath := range oldFilePaths {
		if err := os.Remove(oldFilePath); err != nil {
			log.Warnf("Failed to delete old file: %s", err.Error())
		}
	}

	// **Ambil tugas yang baru di-update beserta file-nya**
	if err := db.Preload("TugasFile").First(&tugas, tugas.ID).Error; err != nil {
		log.WithError(err).Error("Failed to retrieve updated tugas")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to retrieve updated tugas", nil, nil, nil)
	}

	// **Konversi ke Response DTO**
	var tugasResponse dto.TugasResponse
	if err := copier.CopyWithOption(&tugasResponse, &tugas, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{"tugas_id": tugas.ID}).Info("Tugas updated successfully")
	return utils.Response(c, fiber.StatusOK, "Tugas updated successfully", tugasResponse, nil, nil)
}

// DeleteTugas handler untuk menghapus tugas
func DeleteTugas(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "delete_tugas"})

	db := config.DB
	tx := db.Begin() // Mulai transaksi

	tugasID := c.Params("tugasId") // Ambil tugas ID dari parameter
	var tugas models.Tugas

	// Cek apakah tugas ada
	if err := tx.Preload("TugasFile").First(&tugas, tugasID).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("Tugas not found")
		return utils.Response(c, fiber.StatusNotFound, "Tugas not found", nil, nil, nil)
	}

	// Hapus file terkait dari storage
	for _, file := range tugas.TugasFile {
		filePath := fmt.Sprintf("./public/uploads/%s", file.FileURL) // Sesuaikan path
		if err := os.Remove(filePath); err != nil {
			log.Warnf("Failed to delete file: %s", err.Error()) // Log jika gagal hapus
		}
	}

	// Hapus tugas dan otomatis hapus `tugas_files` karena CASCADE
	if err := tx.Delete(&tugas).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("Failed to delete tugas")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to delete tugas", nil, nil, nil)
	}

	tx.Commit() // Commit transaksi jika berhasil
	log.WithFields(logrus.Fields{"tugas_id": tugas.ID}).Info("Tugas deleted successfully")
	return utils.Response(c, fiber.StatusOK, "Tugas deleted successfully", nil, nil, nil)
}
