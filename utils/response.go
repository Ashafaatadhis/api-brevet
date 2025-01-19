package utils

import "github.com/gofiber/fiber/v2"

// Response adalah fungsi untuk menstandarkan response API
func Response(c *fiber.Ctx, status int, message string, data interface{}, user interface{}, token *string) error {
	// Membuat response standar dengan status, message, dan data
	response := fiber.Map{
		"status":  status,
		"message": message,
		"data":    data,
	}

	// Jika user diberikan, tambahkan ke response
	if user != nil {
		response["user"] = user
	}

	if token != nil {
		response["token"] = token
	}

	return c.Status(status).JSON(response)
}
