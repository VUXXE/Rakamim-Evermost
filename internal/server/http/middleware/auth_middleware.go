package middleware

import (
	"strings"

	"evermos-api/internal/helper"
	"evermos-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	LocalsUserID  = "user_id"
	LocalsEmail   = "email"
	LocalsIsAdmin = "is_admin"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "missing Authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid Authorization header format")
		}

		tokenStr := parts[1]
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", err.Error())
		}

		c.Locals(LocalsUserID, claims.UserID)
		c.Locals(LocalsEmail, claims.Email)
		c.Locals(LocalsIsAdmin, claims.IsAdmin)

		return c.Next()
	}
}

func AdminOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals(LocalsIsAdmin).(bool)
		if !ok || !isAdmin {
			return helper.Error(c, fiber.StatusForbidden, "forbidden: only admin can perform this operation", "Forbidden")
		}
		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) uint {
	val := c.Locals(LocalsUserID)
	if id, ok := val.(uint); ok {
		return id
	}
	return 0
}

func GetIsAdmin(c *fiber.Ctx) bool {
	val := c.Locals(LocalsIsAdmin)
	if isAdmin, ok := val.(bool); ok {
		return isAdmin
	}
	return false
}
