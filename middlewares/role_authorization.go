package middlewares

import (
	"new-brevet-be/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// RoleAuthorization fungsi middleware untuk pengecekan role apakah sudah sesuai atau belum sesuai yang dimasukkan di parameter
func RoleAuthorization(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(User)
		log := logrus.WithFields(logrus.Fields{
			"user_id": user.ID,
			"event":   "role_authorization",
		})
		// Periksa apakah role user ada di daftar allowedRoles
		for _, role := range allowedRoles {
			if role == user.Role {
				log.Infof("Role %s Allowed", user.Role)
				return c.Next() // Izinkan akses
			}
		}

		log.Warnf("WARNING: Role %s not allowed", user.Role)
		// Role tidak cocok
		return utils.Response(c, fiber.StatusForbidden, "Access Denied", nil, nil, nil)
	}
}
