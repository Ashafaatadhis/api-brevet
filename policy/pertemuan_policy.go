package policy

import (
	"new-brevet-be/config"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// // GroupBatchOwnerPolicy memeriksa apakah guru hanya bisa mengakses pertemuan berdasarkan group_batches_id miliknya sekaligus validasi exist
// func GroupBatchOwnerPolicy(action string) fiber.Handler {
// 	return func(c *fiber.Ctx) error {

// 		userData := c.Locals("user").(middlewares.User)
// 		grBatchID := c.Params("id")
// 		db := config.DB

// 		// Validasi untuk GET
// 		if action == "view" {
// 			var groupBatch models.GroupBatch
// 			if err := db.Where("id = ? AND teacher_id = ?", grBatchID, userData.ID).First(&groupBatch).Error; err != nil {
// 				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to view this class", nil, nil, nil)
// 			}
// 		}

// 		if action == "delete" {
// 			pertemuanIDStr := c.Params("pertemuanId")
// 			if pertemuanIDStr == "" {
// 				return utils.Response(c, fiber.StatusBadRequest, "Pertemuan ID is required", nil, nil, nil)
// 			}

// 			// Convert pertemuan ID ke integer
// 			pertemuanID, err := strconv.Atoi(pertemuanIDStr)
// 			if err != nil {
// 				return utils.Response(c, fiber.StatusBadRequest, "Invalid Pertemuan ID", nil, nil, nil)
// 			}

// 			// Ambil data Pertemuan dari database
// 			var pertemuan models.Pertemuan
// 			if err := db.Preload("GroupBatch").First(&pertemuan, pertemuanID).Error; err != nil {
// 				log.Println("Error fetching Pertemuan:", err)
// 				return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
// 			}

// 			// Pastikan user yang login adalah guru yang memiliki pertemuan ini
// 			if utils.GetIntValue(pertemuan.GroupBatch.TeacherID) != userData.ID {
// 				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to access this Pertemuan", nil, nil, nil)
// 			}
// 		}

// 		// Untuk delete, kita ambil pertemuan ID dari URL params
// 		if action == "update" {

// 			pertemuanIDStr := c.Params("pertemuanId")
// 			if pertemuanIDStr == "" {
// 				return utils.Response(c, fiber.StatusBadRequest, "Pertemuan ID is required", nil, nil, nil)
// 			}

// 			// Convert pertemuan ID ke integer
// 			pertemuanID, err := strconv.Atoi(pertemuanIDStr)
// 			if err != nil {
// 				return utils.Response(c, fiber.StatusBadRequest, "Invalid Pertemuan ID", nil, nil, nil)
// 			}

// 			// Ambil data Pertemuan dari database
// 			var pertemuan models.Pertemuan
// 			if err := db.Preload("GroupBatch").Preload("GroupBatch.Teacher").First(&pertemuan, pertemuanID).Error; err != nil {
// 				log.Println("Error fetching Pertemuan:", err)
// 				return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
// 			}

// 			if utils.GetIntValue(pertemuan.GroupBatch.TeacherID) != userData.ID {

// 				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to access this Pertemuan", nil, nil, nil)
// 			}

// 			var groupBatch models.GroupBatch
// 			if err := db.First(&groupBatch, grBatchID).Error; err != nil {
// 				log.Println("Error fetching GroupBatch:", err)
// 				return utils.Response(c, fiber.StatusNotFound, "Group batch not found", nil, nil, nil)
// 			}

// 			// Pastikan user yang login adalah guru yang memiliki group batch ini
// 			if utils.GetIntValue(groupBatch.TeacherID) != userData.ID {
// 				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to create this Pertemuan", nil, nil, nil)
// 			}

// 			// Pastikan user yang login adalah guru yang memiliki group batch ini
// 			if utils.GetIntValue(pertemuan.GroupBatch.TeacherID) != userData.ID {
// 				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to create this Pertemuan", nil, nil, nil)
// 			}

// 		}

// 		// Untuk create, kita pastikan group_batches_id milik user yang login
// 		if action == "create" {

// 			var groupBatch models.GroupBatch
// 			if err := db.First(&groupBatch, grBatchID).Error; err != nil {
// 				log.Println("Error fetching GroupBatch:", err)
// 				return utils.Response(c, fiber.StatusNotFound, "Group batch not found", nil, nil, nil)
// 			}

// 			// Pastikan user yang login adalah guru yang memiliki group batch ini
// 			if *groupBatch.TeacherID != userData.ID {
// 				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to create this Pertemuan", nil, nil, nil)
// 			}
// 		}

// 		// Lanjutkan ke handler selanjutnya jika validasi lolos
// 		return c.Next()

// 	}
// }

// GroupBatchOwnerPolicy memeriksa apakah guru hanya bisa mengakses pertemuan berdasarkan group_batches_id miliknya sekaligus validasi exist
func GroupBatchOwnerPolicy(action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userData := c.Locals("user").(middlewares.User)
		grBatchID := c.Params("id")
		db := config.DB

		// Validasi apakah group batch ini dimiliki oleh user yang login
		var groupBatch models.GroupBatch
		if err := db.Where("id = ? AND teacher_id = ?", grBatchID, userData.ID).First(&groupBatch).Error; err != nil {
			return utils.Response(c, fiber.StatusForbidden, "You are not authorized to access this class", nil, nil, nil)
		}

		// Jika hanya GET daftar pertemuan, cukup validasi ownership batch
		if action == "view" {
			return c.Next()
		}

		// Untuk aksi yang berhubungan dengan pertemuan (update, delete, detail pertemuan)
		pertemuanIDStr := c.Params("pertemuanId")
		if pertemuanIDStr != "" {
			pertemuanID, err := strconv.Atoi(pertemuanIDStr)
			if err != nil {
				return utils.Response(c, fiber.StatusBadRequest, "Invalid Pertemuan ID", nil, nil, nil)
			}

			// Validasi apakah pertemuan ini benar-benar ada dan milik guru yang bersangkutan
			var pertemuan models.Pertemuan
			if err := db.Where("id = ? AND gr_batch_id = ?", pertemuanID, grBatchID).First(&pertemuan).Error; err != nil {
				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to access this Pertemuan", nil, nil, nil)
			}

			// Untuk update atau delete, validasi tambahan (hanya guru pemilik yang boleh edit/delete)
			if action == "update" || action == "delete" {
				if *groupBatch.TeacherID != userData.ID {
					return utils.Response(c, fiber.StatusForbidden, "You are not authorized to modify this Pertemuan", nil, nil, nil)
				}
			}
		}

		// Untuk create, pastikan user yang login adalah guru yang memiliki group batch ini
		if action == "create" {
			if *groupBatch.TeacherID != userData.ID {
				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to create a Pertemuan in this class", nil, nil, nil)
			}
		}

		// Lanjut ke handler jika validasi lolos
		return c.Next()
	}
}
