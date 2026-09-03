package utils_test

import (
	"testing"

	"evermos-api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashing(t *testing.T) {
	rawPassword := "SecureP@ssw0rd123!"

	// 1. Hash password
	hash, err := utils.HashPassword(rawPassword)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, rawPassword, hash)

	// 2. Compare correct password
	assert.True(t, utils.CheckPasswordHash(rawPassword, hash))

	// 3. Compare incorrect password
	assert.False(t, utils.CheckPasswordHash("WrongPassword999!", hash))
}
