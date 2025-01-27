package utils

import "github.com/gofiber/fiber/v2"

// ResponseFormat adalah format standar untuk respons API
type responseFormat struct {
	Status  int         `json:"status"`           // Status (true untuk sukses, false untuk error)
	Message string      `json:"message"`          // Pesan untuk client
	Data    interface{} `json:"data,omitempty"`   // Data hasil response (opsional)
	Meta    interface{} `json:"meta,omitempty"`   // Metadata tambahan (opsional)
	Errors  interface{} `json:"errors,omitempty"` // Detail error (opsional)
}

// NewResponse adalah utilitas untuk membuat response API
func NewResponse(c *fiber.Ctx, status int, message string, data interface{}, meta interface{}, errors interface{}) error {
	response := responseFormat{
		Status:  status,
		Message: message,
		Data:    data,
		Meta:    meta,
		Errors:  errors,
	}

	return c.Status(status).JSON(response)
}

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
