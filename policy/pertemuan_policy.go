package policy

import (
	"log"
	"new-brevet-be/config"
	"new-brevet-be/dto"
	"new-brevet-be/middlewares"
	"new-brevet-be/models"
	"new-brevet-be/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GroupBatchOwnerPolicy memeriksa apakah guru hanya bisa mengakses pertemuan berdasarkan group_batches_id miliknya sekaligus validasi exist
func GroupBatchOwnerPolicy(action string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userData := c.Locals("user").(middlewares.User)
		db := config.DB

		if action == "delete" {
			pertemuanIDStr := c.Params("id")
			if pertemuanIDStr == "" {
				return utils.Response(c, fiber.StatusBadRequest, "Pertemuan ID is required", nil, nil, nil)
			}

			// Convert pertemuan ID ke integer
			pertemuanID, err := strconv.Atoi(pertemuanIDStr)
			if err != nil {
				return utils.Response(c, fiber.StatusBadRequest, "Invalid Pertemuan ID", nil, nil, nil)
			}

			// Ambil data Pertemuan dari database
			var pertemuan models.Pertemuan
			if err := db.First(&pertemuan, pertemuanID).Error; err != nil {
				log.Println("Error fetching Pertemuan:", err)
				return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
			}

			// Pastikan user yang login adalah guru yang memiliki pertemuan ini
			if pertemuan.GrBatchID != userData.ID {
				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to access this Pertemuan", nil, nil, nil)
			}
		}

		// Untuk delete, kita ambil pertemuan ID dari URL params
		if action == "update" {

			body := c.Locals("body").(dto.EditPertemuanRequest)
			pertemuanIDStr := c.Params("id")
			if pertemuanIDStr == "" {
				return utils.Response(c, fiber.StatusBadRequest, "Pertemuan ID is required", nil, nil, nil)
			}

			// Convert pertemuan ID ke integer
			pertemuanID, err := strconv.Atoi(pertemuanIDStr)
			if err != nil {
				return utils.Response(c, fiber.StatusBadRequest, "Invalid Pertemuan ID", nil, nil, nil)
			}

			// Ambil data Pertemuan dari database
			var pertemuan models.Pertemuan
			if err := db.Preload("GroupBatch").Preload("GroupBatch.Teacher").First(&pertemuan, pertemuanID).Error; err != nil {
				log.Println("Error fetching Pertemuan:", err)
				return utils.Response(c, fiber.StatusNotFound, "Pertemuan not found", nil, nil, nil)
			}

			if utils.GetIntValue(pertemuan.GroupBatch.TeacherID) != userData.ID {

				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to access this Pertemuan", nil, nil, nil)
			}

			if body.GrBatchID != nil {
				var groupBatch models.GroupBatch
				if err := db.First(&groupBatch, *body.GrBatchID).Error; err != nil {
					log.Println("Error fetching GroupBatch:", err)
					return utils.Response(c, fiber.StatusNotFound, "Group batch not found", nil, nil, nil)
				}

				// Pastikan user yang login adalah guru yang memiliki group batch ini
				if utils.GetIntValue(groupBatch.TeacherID) != userData.ID {
					return utils.Response(c, fiber.StatusForbidden, "You are not authorized to create this Pertemuan", nil, nil, nil)
				}

			}

			// Pastikan user yang login adalah guru yang memiliki group batch ini
			if utils.GetIntValue(pertemuan.GroupBatch.TeacherID) != userData.ID {
				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to create this Pertemuan", nil, nil, nil)
			}

		}

		// Untuk create, kita pastikan group_batches_id milik user yang login
		if action == "create" {
			body := c.Locals("body").(dto.CreatePertemuanRequest)
			var groupBatch models.GroupBatch
			if err := db.First(&groupBatch, body.GrBatchID).Error; err != nil {
				log.Println("Error fetching GroupBatch:", err)
				return utils.Response(c, fiber.StatusNotFound, "Group batch not found", nil, nil, nil)
			}

			// Pastikan user yang login adalah guru yang memiliki group batch ini
			if *groupBatch.TeacherID != userData.ID {
				return utils.Response(c, fiber.StatusForbidden, "You are not authorized to create this Pertemuan", nil, nil, nil)
			}
		}

		// Lanjutkan ke handler selanjutnya jika validasi lolos
		return c.Next()

	}
}
