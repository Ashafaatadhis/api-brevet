package middlewares

import (
	"new-brevet-be/utils"

	"github.com/gofiber/fiber/v2"
)

// RoleAuthorization fungsi middleware untuk pengecekan role apakah sudah sesuai atau belum sesuai yang dimasukkan di parameter
func RoleAuthorization(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(User)

		// Periksa apakah role user ada di daftar allowedRoles
		for _, role := range allowedRoles {
			if role == user.Role {
				return c.Next() // Izinkan akses
			}
		}

		// Role tidak cocok
		return utils.Response(c, fiber.StatusForbidden, "Access Denied", nil, nil, nil)
	}
}
