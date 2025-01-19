package main

import (
	"new-brevet-be/config"
	"new-brevet-be/routes"
	"new-brevet-be/tasks"
	"new-brevet-be/utils"
	"new-brevet-be/validation"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Inisialisasi aplikasi Fiber
	config.LoadEnv()

	validation.InitValidator()

	app := fiber.New()

	// Menyajikan file statis di folder 'public' (misalnya, gambar di folder uploads)
	app.Static("/uploads", "./public/uploads") // File dalam ./public/uploads bisa diakses melalui http://localhost:3000/uploads

	// // Inisialisasi koneksi ke database
	config.InitDB()

	api := app.Group("/api") // /api

	// testing cicd purpose
	api.Get("/hello", func(c *fiber.Ctx) error {
		return utils.Response(c, fiber.StatusOK, "Hello World", nil, nil, nil)
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
