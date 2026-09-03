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
	"evermos-api/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAndStoreHandler_E2E(t *testing.T) {
	app := setupTestApp(t)

	ts := time.Now().UnixNano()
	userAEmail := fmt.Sprintf("e2e_a_%d@evermos.com", ts)
	userAPhone := fmt.Sprintf("0881%d", ts%100000000)

	// Register User A
	regPayloadA := map[string]string{
		"name":     "E2E User A",
		"email":    userAEmail,
		"phone":    userAPhone,
		"password": "Password123!",
	}
	bodyA, _ := json.Marshal(regPayloadA)
	regReqA := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(bodyA))
	regReqA.Header.Set("Content-Type", "application/json")
	regRespA, err := app.Test(regReqA)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, regRespA.StatusCode)

	var resA struct {
		Data model.AuthResponse `json:"data"`
	}
	json.NewDecoder(regRespA.Body).Decode(&resA)
	tokenA := resA.Data.Token
	storeIDA := resA.Data.Store.ID

	// Register User B
	userBEmail := fmt.Sprintf("e2e_b_%d@evermos.com", ts)
	userBPhone := fmt.Sprintf("0882%d", ts%100000000)
	regPayloadB := map[string]string{
		"name":     "E2E User B",
		"email":    userBEmail,
		"phone":    userBPhone,
		"password": "Password123!",
	}
	bodyB, _ := json.Marshal(regPayloadB)
	regReqB := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(bodyB))
	regReqB.Header.Set("Content-Type", "application/json")
	regRespB, err := app.Test(regReqB)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, regRespB.StatusCode)

	var resB struct {
		Data model.AuthResponse `json:"data"`
	}
	json.NewDecoder(regRespB.Body).Decode(&resB)
	tokenB := resB.Data.Token

	// 1. GET /api/v1/users/me with User A
	reqMe := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+tokenA)
	respMe, err := app.Test(reqMe)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respMe.StatusCode)

	// 2. GET /api/v1/stores/me with User A
	reqStoreMe := httptest.NewRequest(http.MethodGet, "/api/v1/stores/me", nil)
	reqStoreMe.Header.Set("Authorization", "Bearer "+tokenA)
	respStoreMe, err := app.Test(reqStoreMe)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respStoreMe.StatusCode)

	// 3. User B tries to update User A's store -> MUST RETURN 403 Forbidden!
	tamperPayload := map[string]string{
		"name": "Tampered Store",
	}
	bodyTamper, _ := json.Marshal(tamperPayload)
	reqTamper := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/stores/%d", storeIDA), bytes.NewReader(bodyTamper))
	reqTamper.Header.Set("Content-Type", "application/json")
	reqTamper.Header.Set("Authorization", "Bearer "+tokenB) // Using User B's token
	respTamper, err := app.Test(reqTamper)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, respTamper.StatusCode)

	// 4. User A legitimately updates User A's store -> MUST SUCCEED (200 OK)
	validPayload := map[string]string{
		"name":    "User A Brand New Store",
		"address": "Jakarta Selatan",
	}
	bodyValid, _ := json.Marshal(validPayload)
	reqValid := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/stores/%d", storeIDA), bytes.NewReader(bodyValid))
	reqValid.Header.Set("Content-Type", "application/json")
	reqValid.Header.Set("Authorization", "Bearer "+tokenA)
	respValid, err := app.Test(reqValid)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respValid.StatusCode)

	// 5. Public GET /api/v1/stores (no auth required)
	reqPublic := httptest.NewRequest(http.MethodGet, "/api/v1/stores?limit=10&offset=0", nil)
	respPublic, err := app.Test(reqPublic)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respPublic.StatusCode)

	var resPublic helper.BaseResponse
	err = json.NewDecoder(respPublic.Body).Decode(&resPublic)
	require.NoError(t, err)
	assert.Equal(t, 200, resPublic.Code)
}
