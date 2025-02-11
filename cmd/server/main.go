package main

import (
	"new-brevet-be/config"
	"new-brevet-be/handlers"
	"new-brevet-be/routes"
	"new-brevet-be/tasks"
	"new-brevet-be/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func main() {
	config.LoadEnv()
	validation.InitValidator()

	config.InitLogger()

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100MB
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://brevet-tax-center.vercel.app",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Custom-Header",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	app.Static("/uploads", "./public/uploads")

	config.InitDB()

	// Middleware global untuk error handling
	app.Use(handlers.ErrorHandler)
	app.Use(helmet.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Asia/Jakarta",
	}))

	app.Get("/metrics", monitor.New())
	api := app.Group("/api")

	api.Get("/hello", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  200,
			"message": "Hello World testing env3",
			"data":    nil,
		})
	})

	v1 := api.Group("/v1")
	routes.Setup(v1)

	go tasks.CleanupExpiredTokens()

	logrus.Info("Application Starting...")

	app.Listen(":3000")
}
