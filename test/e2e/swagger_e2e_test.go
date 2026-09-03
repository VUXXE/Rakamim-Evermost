package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"evermos-api/internal/server/http/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSwaggerHandler(t *testing.T) {
	app := fiber.New()
	swaggerHandler := handler.NewSwaggerHandler()

	app.Get("/swagger", swaggerHandler.GetSwaggerUI)
	app.Get("/swagger/doc.json", swaggerHandler.GetSwaggerJSON)

	// 1. Test Swagger UI HTML
	reqUI := httptest.NewRequest(http.MethodGet, "/swagger", nil)
	respUI, err := app.Test(reqUI)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUI.StatusCode)
	assert.Contains(t, respUI.Header.Get("Content-Type"), "text/html")

	// 2. Test Swagger Spec JSON
	reqJSON := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	respJSON, err := app.Test(reqJSON)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respJSON.StatusCode)
	assert.Contains(t, respJSON.Header.Get("Content-Type"), "application/json")
}
