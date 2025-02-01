package handlers

import (
	"log"
	"new-brevet-be/utils"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler Middleware untuk menangani error global
func ErrorHandler(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v", r)
			_ = utils.Response(c, fiber.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
		}
	}()

	if err := c.Next(); err != nil {
		log.Printf("Error: %v", err)

		var statusCode int
		var message string

		switch e := err.(type) {
		case *fiber.Error:
			statusCode = e.Code
			message = e.Message
		default:
			statusCode = fiber.StatusInternalServerError
			message = "Internal Server Error"
		}

		return utils.Response(c, statusCode, message, nil, nil, nil)
	}

	return nil
}
