package http

import (
	"evermos-api/internal/helper"
	"evermos-api/internal/infrastructure/container"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type RouterConfig struct {
	App       *fiber.App
	Container *container.Container
}

func SetupRoutes(cfg RouterConfig) {
	app := cfg.App
	c := cfg.Container

	// Global Middlewares
	app.Use(recover.New())
	app.Use(logger.New())

	// Serve uploaded files statically
	app.Static("/uploads", "./uploads")

	// Base API Group
	api := app.Group("/api/v1")

	// Health Check
	api.Get("/health", func(ctx *fiber.Ctx) error {
		return helper.Success(ctx, fiber.StatusOK, "evermos e-commerce api is healthy", fiber.Map{
			"status": "online",
		})
	})

	// Auth Endpoints (Public)
	if c != nil && c.AuthHandler != nil {
		auth := api.Group("/auth")
		auth.Post("/register", c.AuthHandler.Register)
		auth.Post("/login", c.AuthHandler.Login)
	}
}
