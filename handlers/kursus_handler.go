package handlers

import (
	"fmt"
	"log"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"new-brevet-be/validation"
	"os"
	"strings"

	dto_mapper "github.com/dranikpg/dto-mapper" // Impor dengan alias
	"github.com/gofiber/fiber/v2"
)

// GetKursus handler untuk mengambil semua kursus dengan preload semua relasi
func GetKursus(c *fiber.Ctx) error {
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
	validSortFields, err := utils.GetValidSortFields(&models.Batch{})
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

	// Mengambil semua kursus dengan preload semua relasi
	var kursusList []models.Kursus
	query := db.Model(&models.Kursus{}).Preload("Teacher").
		Preload("Jenis").
		Preload("GroupBatches").
		Preload("Kelas").
		Preload("Category").
		Preload("Hari")

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

	// Hitung total data sebelum pagination
	var totalData int64
	if err := query.Count(&totalData).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to count total data", nil, nil, err.Error())
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&kursusList).Error; err != nil {
		return utils.NewResponse(c, fiber.StatusInternalServerError, "Failed to get mapping batch", nil, nil, err.Error())
	}

	var kursusResponseList []dto.KursusResponse

	// Automapping
	if err := dto_mapper.Map(&kursusResponseList, kursusList); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map kursus response", nil, nil, nil)
	}

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "Kursus retrieved successfully", kursusResponseList, meta, nil)

}

// GetDetailKursus handler untuk mengambil detail kursus dengan preload semua relasi
func GetDetailKursus(c *fiber.Ctx) error {
	db := config.DB
	kursusID := c.Params("id")

	// Mengambil kursus berdasarkan ID dengan preload semua relasi
	var kursus models.Kursus
	if err := db.Where("id = ?", kursusID).
		Preload("Teacher").
		Preload("Jenis").
		Preload("GroupBatches").
		Preload("Kelas").
		Preload("Category").
		Preload("Hari"). // Preload relasi many-to-many dengan Hari
		First(&kursus).Error; err != nil {
		log.Println("Failed to fetch kursus with relations:", err)
		return utils.Response(c, fiber.StatusNotFound, "Kursus not found", nil, nil, nil)
	}

	// Inisialisasi response
	var kursusResponseList dto.KursusResponse

	// Automapping
	if err := dto_mapper.Map(&kursusResponseList, kursus); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map kursus response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Kursus retrieved successfully", kursusResponseList, nil, nil)
}

// PostKursus adalah handler untuk route post kursus
func PostKursus(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(validation.PostKursus)

	tx := db.Begin()

	kursus := models.Kursus{
		Judul:            body.Judul,
		JenisID:          body.JenisID,
		KelasID:          body.KelasID,
		DeskripsiSingkat: body.DeskripsiSingkat,
		Deskripsi:        body.Deskripsi,
		Pembelajaran:     body.Pembelajaran,
		Diperoleh:        body.Diperoleh,
		CategoryID:       body.CategoryID,
		ThumbnailKursus:  body.ThumbnailKursus,
		ThumbnailURL:     body.ThumbnailURL,
		HargaAsli:        body.HargaAsli,
		HargaDiskon:      body.HargaDiskon,
		StartDate:        body.StartDate,
		EndDate:          body.EndDate,
		StartTime:        body.StartTime,
		EndTime:          body.EndTime,
	}

	if err := tx.Create(&kursus).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create Kursus", nil, nil, nil)
	}

	var hariList []models.Hari
	for _, hariID := range body.HariID {
		var hari models.Hari
		if err := db.First(&hari, hariID).Error; err != nil {
			tx.Rollback()
			return utils.Response(c, fiber.StatusBadRequest, "Invalid 'hari' ID", nil, nil, nil)
		}
		hariList = append(hariList, hari)
	}
	// var hari models.Hari
	// if err := db.First(&hari, body.HariID).Error; err != nil {
	// 	tx.Rollback()
	// 	return utils.Response(c, fiber.StatusBadRequest, "Invalid 'hari' ID", nil, nil, nil)
	// }
	// log.Printf("kursus: %+v", kursus)
	// log.Printf("hariList: %+v", hari)

	if err := tx.Model(&kursus).Association("Hari").Append(&hariList); err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to associate Hari", nil, nil, nil)
	}

	// Commit transaksi
	tx.Commit()

	// upload gambar
	thumbnail, err := c.FormFile("thumbnail_kursus")
	if err != nil {
		log.Print(err.Error())
	}

	var data *string
	path := "thumbnail_kursus"
	if thumbnail != nil {
		data, err = utils.UploadFileHandler(c, thumbnail, &path)
		if err != nil {
			return err
		}
	}
	// Jika gambar berhasil diupload dan transaksi berhasil, update gambar ke dalam body
	if data != nil {

		kursus.ThumbnailKursus = *data
		if err := db.Model(&kursus).Updates(map[string]interface{}{"ThumbnailKursus": kursus.ThumbnailKursus}).Error; err != nil {
			// Jika gagal update, kita bisa rollback atau menangani dengan cara lain
			log.Println("Failed to update kursus with image:", err)
		}
	}

	// Mengambil data kursus dengan preload semua relasi
	var kursusList models.Kursus
	if err := db.Preload("Jenis").
		Preload("Kelas").
		Preload("Category").
		Preload("Hari"). // Preload relasi many-to-many dengan Hari
		First(&kursusList, kursus.ID).Error; err != nil {
		log.Println("Failed to fetch kursus with relations:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get kursus", nil, nil, nil)
	}

	// Inisialisasi response
	var kursusResponseList dto.KursusResponse

	// Automapping
	if err := dto_mapper.Map(&kursusResponseList, kursusList); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Kursus created successfully", kursusResponseList, nil, nil)
}

// UpdateKursus adalah handler untuk route update kursus
func UpdateKursus(c *fiber.Ctx) error {
	db := config.DB
	body := c.Locals("body").(validation.PostKursus)
	kursusID := c.Params("id")

	tx := db.Begin()

	// Cari kursus yang akan diupdate
	var kursus models.Kursus
	if err := db.First(&kursus, kursusID).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusNotFound, "Kursus not found", nil, nil, nil)
	}

	// temp gambar lama
	tempThubmnail := kursus.ThumbnailKursus

	// Update data kursus
	kursus.Judul = body.Judul
	kursus.JenisID = body.JenisID
	kursus.KelasID = body.KelasID
	kursus.DeskripsiSingkat = body.DeskripsiSingkat
	kursus.Deskripsi = body.Deskripsi
	kursus.Pembelajaran = body.Pembelajaran
	kursus.Diperoleh = body.Diperoleh
	kursus.CategoryID = body.CategoryID
	kursus.ThumbnailURL = body.ThumbnailURL
	kursus.HargaAsli = body.HargaAsli
	kursus.HargaDiskon = body.HargaDiskon
	kursus.StartDate = body.StartDate
	kursus.EndDate = body.EndDate
	kursus.StartTime = body.StartTime
	kursus.EndTime = body.EndTime

	if err := tx.Save(&kursus).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to update Kursus", nil, nil, nil)
	}

	// Mengasosiasikan Hari
	var hariList []models.Hari
	for _, hariID := range body.HariID {
		var hari models.Hari
		if err := db.First(&hari, hariID).Error; err != nil {
			tx.Rollback()
			return utils.Response(c, fiber.StatusBadRequest, "Invalid 'hari' ID", nil, nil, nil)
		}
		hariList = append(hariList, hari)
	}

	// Mengupdate asosiasi Hari
	if err := tx.Model(&kursus).Association("Hari").Replace(hariList); err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to associate Hari", nil, nil, nil)
	}

	// Commit transaksi
	tx.Commit()

	// upload gambar jika ada
	thumbnail, err := c.FormFile("thumbnail_kursus")
	if err != nil {
		log.Print(err.Error())
	}

	var data *string
	path := "thumbnail_kursus"
	if thumbnail != nil {
		data, err = utils.UploadFileHandler(c, thumbnail, &path)
		if err != nil {
			return err
		}
	}

	// Jika gambar berhasil diupload, update gambar ke dalam body
	if data != nil {
		kursus.ThumbnailKursus = *data
		if err := db.Model(&kursus).Updates(map[string]interface{}{"ThumbnailKursus": kursus.ThumbnailKursus}).Error; err != nil {
			// Jika gagal update gambar, kita bisa menangani dengan cara lain
			log.Println("Failed to update kursus with image:", err)
		}
	}

	// Hapus gambar lama jika ada
	if tempThubmnail != "" {
		oldAvatarPath := fmt.Sprintf("./public/uploads/%s", tempThubmnail) // Sesuaikan path
		if err := os.Remove(oldAvatarPath); err != nil {
			log.Printf("Failed to delete old avatar: %s", err.Error())
		}
	}

	// Mengambil data kursus dengan preload semua relasi
	var kursusList models.Kursus
	if err := db.Preload("Jenis").
		Preload("Kelas").
		Preload("Category").
		Preload("Hari"). // Preload relasi many-to-many dengan Hari
		First(&kursusList, kursus.ID).Error; err != nil {
		log.Println("Failed to fetch kursus with relations:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get kursus", nil, nil, nil)
	}

	// Inisialisasi response
	var kursusResponseList dto.KursusResponse

	// Automapping
	if err := dto_mapper.Map(&kursusResponseList, kursusList); err != nil {
		log.Println("Error during mapping:", err)
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to map batch response", nil, nil, nil)
	}

	return utils.Response(c, fiber.StatusOK, "Kursus updated successfully", kursusResponseList, nil, nil)
}

// DeleteKursus adalah handler untuk route delete kursus
func DeleteKursus(c *fiber.Ctx) error {
	db := config.DB
	kursusID := c.Params("id")

	tx := db.Begin()

	// Cari kursus yang akan dihapus
	var kursus models.Kursus
	if err := db.First(&kursus, kursusID).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusNotFound, "Kursus not found", nil, nil, nil)
	}

	// Simpan gambar lama untuk dihapus setelah kursus dihapus
	tempThumbnail := kursus.ThumbnailKursus

	// Menghapus asosiasi Hari
	if err := tx.Model(&kursus).Association("Hari").Clear(); err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to remove Hari associations", nil, nil, nil)
	}

	// Hapus kursus
	if err := tx.Delete(&kursus).Error; err != nil {
		tx.Rollback()
		return utils.Response(c, fiber.StatusBadRequest, "Failed to delete Kursus", nil, nil, nil)
	}

	// Commit transaksi
	tx.Commit()

	// Hapus gambar terkait jika ada
	if tempThumbnail != "" {
		oldAvatarPath := fmt.Sprintf("./public/uploads/%s", tempThumbnail) // Sesuaikan path
		if err := os.Remove(oldAvatarPath); err != nil {
			log.Printf("Failed to delete old thumbnail: %s", err.Error())
		}
	}

	return utils.Response(c, fiber.StatusOK, "Kursus deleted successfully", nil, nil, nil)
}
