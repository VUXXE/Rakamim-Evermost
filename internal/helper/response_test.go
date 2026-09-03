package helper_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"evermos-api/internal/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseHelpers(t *testing.T) {
	app := fiber.New()

	app.Get("/test/success", func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.StatusOK, "operation successful", map[string]string{"key": "value"})
	})

	app.Get("/test/error", func(c *fiber.Ctx) error {
		return helper.Error(c, fiber.StatusBadRequest, "operation failed", "detailed error message")
	})

	app.Get("/test/pagination", func(c *fiber.Ctx) error {
		items := []string{"item1", "item2"}
		return helper.Pagination(c, fiber.StatusOK, "list fetched", 42, 10, 0, items)
	})

	// 1. Test Success
	reqSuccess := httptest.NewRequest(http.MethodGet, "/test/success", nil)
	respSuccess, err := app.Test(reqSuccess)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respSuccess.StatusCode)

	var successBody helper.BaseResponse
	err = json.NewDecoder(respSuccess.Body).Decode(&successBody)
	require.NoError(t, err)
	assert.Equal(t, 200, successBody.Code)
	assert.Equal(t, "operation successful", successBody.Message)

	// 2. Test Error
	reqError := httptest.NewRequest(http.MethodGet, "/test/error", nil)
	respError, err := app.Test(reqError)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, respError.StatusCode)

	var errorBody helper.BaseResponse
	err = json.NewDecoder(respError.Body).Decode(&errorBody)
	require.NoError(t, err)
	assert.Equal(t, 400, errorBody.Code)
	assert.Equal(t, "operation failed", errorBody.Message)
	assert.Equal(t, "detailed error message", errorBody.Error)

	// 3. Test Pagination
	reqPagination := httptest.NewRequest(http.MethodGet, "/test/pagination", nil)
	respPagination, err := app.Test(reqPagination)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respPagination.StatusCode)

	var pagBody struct {
		Code    int                   `json:"code"`
		Message string                `json:"message"`
		Data    helper.PaginationData `json:"data"`
	}
	err = json.NewDecoder(respPagination.Body).Decode(&pagBody)
	require.NoError(t, err)
	assert.Equal(t, 200, pagBody.Code)
	assert.Equal(t, int64(42), pagBody.Data.Total)
	assert.Equal(t, 10, pagBody.Data.Limit)
	assert.Equal(t, 0, pagBody.Data.Offset)
}
