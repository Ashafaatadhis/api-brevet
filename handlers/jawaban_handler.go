package handlers

import (
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

// SubmitJawaban handler untuk mensubmit jawaban untuk role siswa
func SubmitJawaban(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "submit_jawaban"})

	db := config.DB

	body := c.Locals("body").(dto.SubmitJawabanRequest)

	log.Print("ini body", *body.Answer, len(body.Files))
	// Ambil tugasID dari parameter URL
	tugasIDStr := c.Params("tugasId")
	tugasID, err := strconv.Atoi(tugasIDStr)
	if err != nil {
		log.WithError(err).Error("Invalid tugas_id format")
		return utils.Response(c, fiber.StatusBadRequest, "Invalid tugas ID", nil, nil, nil)
	}

	files, err := utils.ParseMultipartForm(c, "files")
	if err != nil {
		return utils.Response(c, fiber.StatusBadRequest, "Failed to parse form", nil, nil, nil)
	}

	// ✅ Cek apakah dua-duanya kosong
	if (body.Answer == nil || *body.Answer == "") && len(files) == 0 {
		log.Print(len(body.Files), "test isi body")
		return utils.Response(c, fiber.StatusBadRequest, "Jawaban tidak boleh kosong", nil, nil, nil)
	}

	// Start Transaction
	tx := db.Begin()

	// Cek apakah user sudah pernah mengumpulkan jawaban untuk tugas ini
	var existingJawaban models.Jawaban
	if err := tx.Where("tugas_id = ? AND user_id = ?", tugasID, token.ID).First(&existingJawaban).Error; err == nil {
		log.Error("User already submitted an answer")
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "You have already submitted an answer for this task", nil, nil, nil)
	}

	// ✅ Fetch tugas dari database untuk ambil `deadline`
	var tugas models.Tugas
	if err := db.First(&tugas, tugasID).Error; err != nil {
		return utils.Response(c, fiber.StatusNotFound, "Tugas not found", nil, nil, nil)
	}

	var jawaban models.Jawaban

	if err := copier.CopyWithOption(&jawaban, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		tx.Rollback()
		return err
	}

	jawaban.TugasID = tugasID
	jawaban.UserID = token.ID
	jawaban.IsLate = time.Now().After(tugas.EndAt)

	if err := tx.Create(&jawaban).Error; err != nil {
		log.WithError(err).Error("Failed to submit jawaban")
		tx.Rollback()
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to submit jawaban", nil, nil, nil)
	}

	// **Handle Multiple File Upload (Jika Ada)**
	fileURLs, err := utils.UploadMultipleFiles(c, "files", "jawaban")
	if err != nil {
		log.WithError(err).Error("Failed to upload files")
		tx.Rollback()
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to upload files", nil, nil, nil)
	}

	// Simpan file jawaban ke database dengan Association Append
	var jawabanFiles []models.JawabanFile
	for _, fileURL := range fileURLs {
		jawabanFiles = append(jawabanFiles, models.JawabanFile{
			FileURL: fileURL,
		})
	}

	if err := tx.Model(&jawaban).Association("JawabanFile").Append(&jawabanFiles); err != nil {
		tx.Rollback()
		log.WithError(err).Error("Failed to associate files")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to associate files", nil, nil, nil)
	}

	// Commit transaction jika semua berhasil
	tx.Commit()

	// Ambil jawaban yang baru dibuat dengan Preload jawaban_files
	if err := db.Preload("JawabanFile").First(&jawaban, jawaban.ID).Error; err != nil {
		log.WithError(err).Error("Failed to retrieve submitted jawaban")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to retrieve submitted jawaban", nil, nil, nil)
	}

	// Konversi ke Response DTO
	var jawabanResponse dto.JawabanResponse
	if err := copier.CopyWithOption(&jawabanResponse, &jawaban, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{"jawaban_id": jawaban.ID}).Info("Jawaban submitted successfully")
	return utils.Response(c, fiber.StatusOK, "Jawaban submitted successfully", jawabanResponse, nil, nil)
}
