package http

import (
	"evermos-api/internal/helper"
	"evermos-api/internal/infrastructure/container"
	"evermos-api/internal/server/http/middleware"
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

	if c == nil {
		return
	}

	// Auth Endpoints (Public)
	if c.AuthHandler != nil {
		auth := api.Group("/auth")
		auth.Post("/register", c.AuthHandler.Register)
		auth.Post("/login", c.AuthHandler.Login)
	}

	// User Endpoints
	if c.UserHandler != nil {
		users := api.Group("/users", middleware.JWTMiddleware())
		users.Get("/me", c.UserHandler.GetMe)
		users.Patch("/me", c.UserHandler.UpdateMe)
		users.Delete("/me", c.UserHandler.DeleteMe)
		users.Get("/", middleware.AdminOnlyMiddleware(), c.UserHandler.GetAllUsers)
		users.Get("/:id", c.UserHandler.GetUserByID)
	}

	// Store Endpoints
	if c.StoreHandler != nil {
		stores := api.Group("/stores")

		// Protected store routes
		stores.Get("/me", middleware.JWTMiddleware(), c.StoreHandler.GetMyStore)
		stores.Post("/", middleware.JWTMiddleware(), c.StoreHandler.CreateStore)
		stores.Patch("/:id", middleware.JWTMiddleware(), c.StoreHandler.UpdateStore)
		stores.Delete("/:id", middleware.JWTMiddleware(), c.StoreHandler.DeleteStore)

		// Public store browsing routes
		stores.Get("/", c.StoreHandler.GetAllStores)
		stores.Get("/:id", c.StoreHandler.GetStoreByID)
	}

	// Address Endpoints (Protected: Bearer Token Required)
	if c.AddressHandler != nil {
		addresses := api.Group("/addresses", middleware.JWTMiddleware())
		addresses.Post("/", c.AddressHandler.CreateAddress)
		addresses.Get("/", c.AddressHandler.GetMyAddresses)
		addresses.Get("/:id", c.AddressHandler.GetAddressByID)
		addresses.Patch("/:id", c.AddressHandler.UpdateAddress)
		addresses.Delete("/:id", c.AddressHandler.DeleteAddress)
	}
}
