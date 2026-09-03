package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"evermos-api/internal/helper"
	"evermos-api/internal/infrastructure/container"
	"evermos-api/internal/infrastructure/mysql"
	httpserver "evermos-api/internal/server/http"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) *fiber.App {
	cfg := mysql.Config{
		Host:            "127.0.0.1",
		Port:            "3306",
		User:            "root",
		Password:        "password",
		DBName:          "evermos",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	db, err := mysql.Connect(cfg)
	require.NoError(t, err)

	err = mysql.AutoMigrate(db)
	require.NoError(t, err)

	c := container.SetupContainer(db)

	app := fiber.New()
	httpserver.SetupRoutes(httpserver.RouterConfig{
		App:       app,
		Container: c,
	})

	return app
}

func TestAuthHandler_Register_HTTP(t *testing.T) {
	app := setupTestApp(t)

	ts := time.Now().UnixNano()
	payload := map[string]string{
		"name":     fmt.Sprintf("HTTP User %d", ts%10000),
		"email":    fmt.Sprintf("http_%d@evermos.com", ts),
		"phone":    fmt.Sprintf("0852%d", ts%100000000),
		"password": "Password123!",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var res helper.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)
	assert.Equal(t, 201, res.Code)
	assert.Equal(t, "register success", res.Message)
	assert.NotNil(t, res.Data)
}

func TestAuthHandler_Login_HTTP(t *testing.T) {
	app := setupTestApp(t)

	ts := time.Now().UnixNano()
	email := fmt.Sprintf("http_login_%d@evermos.com", ts)
	phone := fmt.Sprintf("0853%d", ts%100000000)
	password := "LoginSecret123!"

	// Register first
	regPayload := map[string]string{
		"name":     "Login User",
		"email":    email,
		"phone":    phone,
		"password": password,
	}
	regBody, _ := json.Marshal(regPayload)
	regReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	_, err := app.Test(regReq)
	require.NoError(t, err)

	// Login
	loginPayload := map[string]string{
		"email":    email,
		"password": password,
	}
	loginBody, _ := json.Marshal(loginPayload)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(loginReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var res helper.BaseResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "login success", res.Message)
}
