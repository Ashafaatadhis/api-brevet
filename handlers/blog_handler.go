package handlers

import (
	"fmt"
	"math"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

// GetAllBlog handler untuk method GET
func GetAllBlog(c *fiber.Ctx) error {
	db := config.DB
	log := logrus.WithFields(logrus.Fields{"event": "get_all_blog"})

	// Ambil query parameters
	search := c.Query("q", "")
	sort := c.Query("sort", "id")
	order := c.Query("order", "asc")
	selectFields := c.Query("select", "")
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	filters := c.Queries() // Mengambil semua query parameter

	// Pagination offset
	offset := (page - 1) * limit

	log.Info("Fetching Blog")

	// Ambil valid sort fields dari tabel blog
	validSortFields, err := utils.GetValidSortFields(&models.Blog{})
	if err != nil {
		log.WithError(err).Error("Failed to get valid sort fields")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to get valid sort fields", nil, nil, nil)
	}

	// Validasi sorting field
	if !validSortFields[sort] {
		sort = "id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	var blogList []models.Blog
	query := db.Model(&models.Blog{})
	// Apply search query
	if search != "" {
		query = query.Where("judul LIKE ?", "%"+search+"%")
	}

	// multi filter
	// Loop semua parameter query
	for field, value := range filters {

		// Cek apakah field valid
		if !validSortFields[field] {
			continue
		}

		// Jika ada koma, gunakan IN query
		values := strings.Split(value, ",")
		if len(values) > 1 {
			query = query.Where(field+" IN ?", values)
		} else {
			query = query.Where(field+" = ?", value)
		}
	}

	// Apply select fields
	if selectFields != "" {
		fields := strings.Split(selectFields, ",")
		query = query.Select(fields)
	}

	// Apply sorting
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	// Hitung total data sebelum pagination
	var totalData int64
	if err := query.Count(&totalData).Error; err != nil {
		log.WithError(err).Error("Failed to count total blog")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to count total blog", nil, nil, nil)
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&blogList).Error; err != nil {
		log.WithError(err).Error("Failed to fetch blog")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to fetch blog", nil, nil, nil)
	}

	log.WithField("total_data", totalData).Info("Successfully fetched blog")

	// Mapping ke DTO response
	var responseList []dto.BlogResponse
	if err := copier.CopyWithOption(&responseList, &blogList, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map blog response")
		return err
	}

	// Metadata pagination
	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total_data": totalData,
		"total_page": int(math.Ceil(float64(totalData) / float64(limit))),
	}

	// Success response
	return utils.NewResponse(c, fiber.StatusOK, "blog fetched successfully", responseList, meta, nil)
}

// GetBlogBySlug handler untuk method GET
func GetBlogBySlug(c *fiber.Ctx) error {
	log := logrus.WithFields(logrus.Fields{"event": "get_blog_by_id"})

	db := config.DB

	slug := c.Params("slug")

	var blog models.Blog

	// Cari blog berdasarkan ID dan preload relasi
	if err := db.Where("slug = ?", slug).
		First(&blog).Error; err != nil {
		log.WithFields(logrus.Fields{"slug": slug}).WithError(err).Error("Failed to fetch blog by ID")
		return utils.Response(c, fiber.StatusNotFound, "Blog not found", nil, nil, nil)
	}

	// Inisialisasi response
	var blogResponse dto.BlogResponse
	if err := copier.CopyWithOption(&blogResponse, &blog, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map blog response")
		return err
	}

	// Set struct kosong menjadi nil
	utils.TransformResponse(&blogResponse)

	log.WithFields(logrus.Fields{"id": blog}).Info("Successfully fetched blog by ID")
	return utils.Response(c, fiber.StatusOK, "Blog fetched successfully", blogResponse, nil, nil)
}

// CreateBlog handler untuk POST blog
func CreateBlog(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "create_blog"})

	db := config.DB
	body := c.Locals("body").(dto.CreateBlogRequest)

	var blog models.Blog

	if err := copier.CopyWithOption(&blog, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map blog response")
		return err
	}

	fileURL, err := utils.UploadImage(c, "gambar", "blog")

	if err != nil {
		log.WithError(err).Error("Failed Upload File")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed Upload File", nil, nil, nil)
	} else if fileURL != nil {
		blog.Gambar = *fileURL
	}

	// Insert ke database
	if err := db.Create(&blog).Scan(&blog).Error; err != nil {
		log.WithError(err).Error("Error creating blog")
		return utils.Response(c, fiber.StatusBadRequest, "Failed to create blog", nil, nil, nil)
	}

	// Automapping response
	var blogResponse dto.BlogResponse
	if err := copier.CopyWithOption(&blogResponse, &blog, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map blog response")
		return err
	}

	log.WithFields(logrus.Fields{"id": blog.ID}).Info("blog created successfully")
	return utils.Response(c, fiber.StatusOK, "Blog created successfully", blogResponse, nil, nil)
}

// UpdateBlog handler untuk PUT blog
func UpdateBlog(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "update_blog"})

	db := config.DB
	body := c.Locals("body").(dto.UpdateBlogRequest)
	slug := c.Params("slug")

	var blog models.Blog

	if err := db.Where("slug = ?", slug).First(&blog).Error; err != nil {
		log.WithError(err).Error("Blog not found")
		return utils.Response(c, fiber.StatusNotFound, "Blog not found", nil, nil, nil)
	}

	if err := copier.CopyWithOption(&blog, &body, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map blog response")
		return err
	}

	// Simpan path file lama untuk nanti dihapus jika ada file baru
	oldFileURL := blog.Gambar

	file, err := c.FormFile("gambar")
	fileURL, err := utils.UploadImage(c, "gambar", "blog")

	if file != nil { // Jika ada file yang dikirim dalam request

		if err != nil {
			log.WithError(err).Error("Failed Upload File")
			return utils.Response(c, fiber.StatusInternalServerError, "Failed Upload File", nil, nil, nil)
		} else if fileURL != nil {
			blog.Gambar = *fileURL
		}
	}

	// Update hanya field yang dikirim dalam request berdasarkan slug
	if err := db.Model(&blog).Updates(blog).Error; err != nil {
		log.WithError(err).Error("Failed to update blog")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to update blog", nil, nil, nil)
	}

	// Hapus gambar lama jika ada
	if oldFileURL != "" {
		oldAvatarPath := fmt.Sprintf("./public/uploads/%s", oldFileURL) // Sesuaikan path
		if err := os.Remove(oldAvatarPath); err != nil {
			log.Warnf("Failed to delete old blog image: %s", err.Error())
		}
	}

	// Automapping response
	var blogResponse dto.BlogResponse
	if err := copier.CopyWithOption(&blogResponse, &blog, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}); err != nil {
		log.WithError(err).Error("Failed to map blog response")
		return err
	}

	log.WithFields(logrus.Fields{"id": blog.ID}).Info("blog created successfully")
	return utils.Response(c, fiber.StatusOK, "Blog created successfully", blogResponse, nil, nil)
}

// DeleteBlog handler untuk Delete Blog
func DeleteBlog(c *fiber.Ctx) error {
	token := c.Locals("user").(middlewares.User)
	log := logrus.WithFields(logrus.Fields{"user_id": token.ID, "event": "delete_blog"})

	db := config.DB

	// Ambil ID Materi dari parameter URL
	slug := c.Params("slug")

	// Cek apakah materi ada
	var blog models.Blog
	if err := db.Where("slug = ?", slug).First(&blog).Error; err != nil {
		log.WithError(err).Error("Blog not found")
		return utils.Response(c, fiber.StatusNotFound, "Blog not found", nil, nil, nil)
	}

	// **Hapus file yang terkait (jika ada)**
	if blog.Gambar != "" {
		oldAvatarPath := fmt.Sprintf("./public/uploads/%s", blog.Gambar) // Sesuaikan path
		if err := os.Remove(oldAvatarPath); err != nil {
			log.Warnf("Failed to delete old blog image: %s", err.Error())
		}
	}

	// **Hapus materi dari database**
	if err := db.Delete(&blog).Error; err != nil {
		log.WithError(err).Error("Failed to delete blog")
		return utils.Response(c, fiber.StatusInternalServerError, "Failed to delete blog", nil, nil, nil)
	}

	log.WithFields(logrus.Fields{"blog_slug": blog.Slug}).Info("Blog deleted successfully")
	return utils.Response(c, fiber.StatusOK, "Blog deleted successfully", nil, nil, nil)
}
