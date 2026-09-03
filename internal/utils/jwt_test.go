package utils_test

import (
	"os"
	"testing"
	"time"

	"evermos-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWT_Lifecycle(t *testing.T) {
	_ = os.Setenv("JWT_SECRET", "test_super_secret_jwt_key_32_bytes_long!")

	userID := uint(42)
	email := "user42@evermos.com"
	isAdmin := true

	// 1. Generate Token
	token, err := utils.GenerateToken(userID, email, isAdmin)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// 2. Validate Token
	claims, err := utils.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, isAdmin, claims.IsAdmin)

	// 3. Validate Invalid Token
	_, err = utils.ValidateToken("invalid.jwt.token")
	assert.Error(t, err)

	// 4. Validate Expired Token
	expiredClaims := utils.JWTClaims{
		UserID:  userID,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	expiredTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenStr, err := expiredTokenObj.SignedString([]byte("test_super_secret_jwt_key_32_bytes_long!"))
	require.NoError(t, err)

	_, err = utils.ValidateToken(expiredTokenStr)
	assert.Error(t, err)
}
