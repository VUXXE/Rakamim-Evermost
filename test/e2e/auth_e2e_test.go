package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"evermos-api/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	resp, err := testReq(app, req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var res struct {
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    model.AuthResponse `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)
	assert.Equal(t, 201, res.Code)
	assert.NotEmpty(t, res.Data.Token)
	assert.Equal(t, fmt.Sprintf("HTTP User %d's Store", ts%10000), res.Data.Store.Name)
}

func TestAuthHandler_Login_HTTP(t *testing.T) {
	app := setupTestApp(t)

	ts := time.Now().UnixNano()
	email := fmt.Sprintf("login_%d@evermos.com", ts)
	phone := fmt.Sprintf("0853%d", ts%100000000)
	password := "Password123!"

	// First register
	regPayload := map[string]string{
		"name":     "Login Tester",
		"email":    email,
		"phone":    phone,
		"password": password,
	}
	regBody, _ := json.Marshal(regPayload)
	regReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	_, err := testReq(app, regReq)
	require.NoError(t, err)

	// Now login
	loginPayload := map[string]string{
		"email":    email,
		"password": password,
	}
	loginBody, _ := json.Marshal(loginPayload)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")

	resp, err := testReq(app, loginReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var res struct {
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    model.AuthResponse `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)
	assert.Equal(t, 200, res.Code)
	assert.NotEmpty(t, res.Data.Token)
}
