package main

import (
	"new-brevet-be/config"

	"new-brevet-be/routes"
	"new-brevet-be/tasks"
	"new-brevet-be/utils"
	"new-brevet-be/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Inisialisasi aplikasi Fiber
	config.LoadEnv()

	validation.InitValidator()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Custom-Header",
	}))

	app.Static("/uploads", "./public/uploads") // File dalam ./public/uploads bisa diakses melalui http://localhost:3000/uploads

	// // Inisialisasi koneksi ke database
	config.InitDB()

	api := app.Group("/api") // /api

	api.Get("/hello", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  200,
			"message": "Hello World jidankuhh",
			"data":    nil,
		})
	})

	v1 := api.Group("/v1") // /api/v1

	// Setup routes// Middleware untuk menangani error
	app.Use(func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			// Tangani error dan kembalikan respons yang sesuai
			var statusCode int
			var message string

			// Ambil error yang dilemparkan
			switch e := err.(type) {
			case *fiber.Error:
				statusCode = e.Code
				message = e.Message
			default:
				statusCode = fiber.StatusInternalServerError
				message = "Internal Server Error"
			}

			// Kirimkan error ke client
			return utils.Response(c, statusCode, message, nil, nil, nil)

		}
		return nil
	})

	routes.Setup(v1)

	go tasks.CleanupExpiredTokens()

	// Menjalankan server
	app.Listen(":3000")
}
