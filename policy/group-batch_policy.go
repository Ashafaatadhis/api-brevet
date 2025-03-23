package policy

import (
	"log"
	"new-brevet-be/config"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"

	"github.com/gofiber/fiber/v2"
)

// GroupBatchAccessPolicy memeriksa apakah user telah membeli kursus sebelum bisa mengakses pertemuan
func GroupBatchAccessPolicy() fiber.Handler {
	return func(c *fiber.Ctx) error {

		userData := c.Locals("user").(middlewares.User)
		grBatchID := c.Params("id") // ID GroupBatch dari URL params
		db := config.DB

		// Cek apakah user sudah membeli kursus dengan GroupBatch ini
		var purchase models.Purchase
		if err := db.Where("user_id = ? AND gr_batch_id = ? AND status_payment_id = ?", userData.ID, grBatchID, 2).First(&purchase).Error; err != nil {
			log.Println("User has not purchased this course:", err)
			return utils.Response(c, fiber.StatusForbidden, "You are not authorized to access this course", nil, nil, nil)
		}

		// Lanjutkan ke handler berikutnya jika lolos validasi
		return c.Next()
	}
}
