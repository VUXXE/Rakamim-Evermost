package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"evermos-api/internal/helper"
	"evermos-api/internal/server/http/middleware"
	"evermos-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddlewares(t *testing.T) {
	_ = os.Setenv("JWT_SECRET", "test_super_secret_jwt_key_32_bytes_long!")

	app := fiber.New()

	// Protected route
	app.Get("/test/protected", middleware.JWTMiddleware(), func(c *fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		isAdmin := middleware.GetIsAdmin(c)
		return helper.Success(c, fiber.StatusOK, "authenticated", fiber.Map{
			"user_id":  userID,
			"is_admin": isAdmin,
		})
	})

	// Admin-only route
	app.Get("/test/admin", middleware.JWTMiddleware(), middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.StatusOK, "admin granted", nil)
	})

	// 1. Missing Authorization Header -> 401 Unauthorized
	reqNoAuth := httptest.NewRequest(http.MethodGet, "/test/protected", nil)
	respNoAuth, err := app.Test(reqNoAuth)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, respNoAuth.StatusCode)

	// 2. Invalid Token Format -> 401 Unauthorized
	reqBadFormat := httptest.NewRequest(http.MethodGet, "/test/protected", nil)
	reqBadFormat.Header.Set("Authorization", "InvalidHeaderFormat")
	respBadFormat, err := app.Test(reqBadFormat)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, respBadFormat.StatusCode)

	// 3. Valid User Token -> 200 OK for protected route
	userToken, err := utils.GenerateToken(10, "user@test.com", false)
	require.NoError(t, err)

	reqUser := httptest.NewRequest(http.MethodGet, "/test/protected", nil)
	reqUser.Header.Set("Authorization", "Bearer "+userToken)
	respUser, err := app.Test(reqUser)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUser.StatusCode)

	// 4. Non-admin accessing Admin route -> 403 Forbidden
	reqForbidden := httptest.NewRequest(http.MethodGet, "/test/admin", nil)
	reqForbidden.Header.Set("Authorization", "Bearer "+userToken)
	respForbidden, err := app.Test(reqForbidden)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, respForbidden.StatusCode)

	// 5. Admin accessing Admin route -> 200 OK
	adminToken, err := utils.GenerateToken(1, "admin@test.com", true)
	require.NoError(t, err)

	reqAdmin := httptest.NewRequest(http.MethodGet, "/test/admin", nil)
	reqAdmin.Header.Set("Authorization", "Bearer "+adminToken)
	respAdmin, err := app.Test(reqAdmin)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respAdmin.StatusCode)
}
