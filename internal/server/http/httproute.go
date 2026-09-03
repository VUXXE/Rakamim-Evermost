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

	// Interactive Swagger UI Documentation
	if c.SwaggerHandler != nil {
		app.Get("/swagger", c.SwaggerHandler.GetSwaggerUI)
		app.Get("/swagger/doc.json", c.SwaggerHandler.GetSwaggerJSON)
		app.Get("/docs", func(ctx *fiber.Ctx) error {
			return ctx.Redirect("/swagger")
		})
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

	// Category Endpoints
	if c.CategoryHandler != nil {
		categories := api.Group("/categories")

		// Admin-Only Category Management Routes
		categories.Post("/", middleware.JWTMiddleware(), middleware.AdminOnlyMiddleware(), c.CategoryHandler.CreateCategory)
		categories.Patch("/:id", middleware.JWTMiddleware(), middleware.AdminOnlyMiddleware(), c.CategoryHandler.UpdateCategory)
		categories.Delete("/:id", middleware.JWTMiddleware(), middleware.AdminOnlyMiddleware(), c.CategoryHandler.DeleteCategory)

		// Public Category Browsing Routes
		categories.Get("/", c.CategoryHandler.GetAllCategories)
		categories.Get("/:id", c.CategoryHandler.GetCategoryByID)
	}

	// Product Endpoints
	if c.ProductHandler != nil {
		products := api.Group("/products")

		// Protected product management routes
		products.Get("/me", middleware.JWTMiddleware(), c.ProductHandler.GetMyProducts)
		products.Post("/", middleware.JWTMiddleware(), c.ProductHandler.CreateProduct)
		products.Patch("/:id", middleware.JWTMiddleware(), c.ProductHandler.UpdateProduct)
		products.Delete("/:id", middleware.JWTMiddleware(), c.ProductHandler.DeleteProduct)

		// Public product browsing and search routes
		products.Get("/", c.ProductHandler.GetAllProducts)
		products.Get("/:id", c.ProductHandler.GetProductByID)
	}

	// Transaction Endpoints (Protected)
	if c.TransactionHandler != nil {
		transactions := api.Group("/transactions", middleware.JWTMiddleware())

		// Buyer personal transaction endpoints
		transactions.Post("/", c.TransactionHandler.CreateTransaction)
		transactions.Get("/me", c.TransactionHandler.GetMyTransactions)
		transactions.Get("/me/:id", c.TransactionHandler.GetMyTransactionByID)
		transactions.Patch("/me/:id", c.TransactionHandler.UpdateTransactionStatus)

		// Admin system-wide transaction endpoints
		transactions.Get("/", middleware.AdminOnlyMiddleware(), c.TransactionHandler.GetAllTransactions)
		transactions.Patch("/:id", middleware.AdminOnlyMiddleware(), c.TransactionHandler.UpdateTransactionStatus)
	}

	// Indonesian Administrative Regions (EMSifa External API Integration)
	if c.RegionHandler != nil {
		regions := api.Group("/regions")
		regions.Get("/provinces", c.RegionHandler.GetProvinces)
		regions.Get("/regencies/:province_id", c.RegionHandler.GetRegencies)
		regions.Get("/districts/:regency_id", c.RegionHandler.GetDistricts)
		regions.Get("/villages/:district_id", c.RegionHandler.GetVillages)
	}
}
