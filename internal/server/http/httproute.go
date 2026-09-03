package http

import (
	"evermos-api/internal/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type RouterConfig struct {
	App *fiber.App
}

func SetupRoutes(cfg RouterConfig) {
	app := cfg.App

	// Global Middlewares
	app.Use(recover.New())
	app.Use(logger.New())

	// Serve uploaded files statically
	app.Static("/uploads", "./uploads")

	// Base API Group
	api := app.Group("/api/v1")

	// Health Check
	api.Get("/health", func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.StatusOK, "evermos e-commerce api is healthy", fiber.Map{
			"status": "online",
		})
	})
}
