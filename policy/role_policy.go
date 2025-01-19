package policy

import (
	"new-brevet-be/config"
	"new-brevet-be/middlewares"
	"new-brevet-be/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type getUserRole struct {
	ID     int `json:"id"`
	RoleID int `json:"role_id"`
}

func (getUserRole) TableName() string {
	return "users"
}

// RolePolicy untuk menentukan hirari level sesuai role
func RolePolicy(c *fiber.Ctx) error {
	db := config.DB
	token := c.Locals("user").(middlewares.User)
	id := c.Params("id")

	var userRole getUserRole

	if err := db.Where("id = ?", id).First(&userRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Response(c, fiber.StatusBadRequest, "user does not exist", nil, nil, nil)
		}
		return utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
	}

	// Pastikan level role target lebih rendah daripada user yang sedang login
	if token.Level >= userRole.RoleID {
		return utils.Response(c, fiber.StatusForbidden, "You are not allowed to manage this role", nil, nil, nil)
	}

	// Lanjutkan ke handler jika validasi lolos
	return c.Next()

}
