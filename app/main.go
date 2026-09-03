package main

import (
	"fmt"
	"log"
	"os"

	"evermos-api/internal/infrastructure/container"
	"evermos-api/internal/infrastructure/mysql"
	httpserver "evermos-api/internal/server/http"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables if .env exists
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, reading from environment variables")
	}

	var appContainer *container.Container
	// Connect to Database
	dbConfig := mysql.LoadConfigFromEnv()
	db, err := mysql.Connect(dbConfig)
	if err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
	} else {
		log.Println("Database connection established successfully")
		if err := mysql.AutoMigrate(db); err != nil {
			log.Fatalf("Auto-migration failed: %v", err)
		}
		appContainer = container.SetupContainer(db)
	}

	// Initialize Fiber App (10MB body limit for image uploads)
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10 MB
		AppName:   "Evermos E-Commerce API",
	})

	// Setup Routes
	httpserver.SetupRoutes(httpserver.RouterConfig{
		App:       app,
		Container: appContainer,
	})

	port := os.Getenv("APP_HTTPPORT")
	if port == "" {
		port = "8080"
	}
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting server on %s ...", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
