package handler

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
)

//go:embed swagger.html
var swaggerHTML string

//go:embed doc.json
var swaggerJSON []byte

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

func (h *SwaggerHandler) GetSwaggerUI(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(swaggerHTML)
}

func (h *SwaggerHandler) GetSwaggerJSON(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.Send(swaggerJSON)
}
