package e2e

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
	"evermos-api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryHandler_AdminRBAC_E2E(t *testing.T) {
	app := setupTestApp(t)

	ts := time.Now().UnixNano()
	normalUserEmail := fmt.Sprintf("normal_%d@evermos.com", ts)
	normalUserPhone := fmt.Sprintf("0896%d", ts%100000000)

	// 1. Register Normal User (is_admin = false)
	regPayload := map[string]string{
		"name":     "Normal User",
		"email":    normalUserEmail,
		"phone":    normalUserPhone,
		"password": "Password123!",
	}
	body, _ := json.Marshal(regPayload)
	regReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	regReq.Header.Set("Content-Type", "application/json")
	regResp, err := testReq(app, regReq)
	require.NoError(t, err)

	var regRes struct {
		Data model.AuthResponse `json:"data"`
	}
	json.NewDecoder(regResp.Body).Decode(&regRes)
	normalToken := regRes.Data.Token
	userID := regRes.Data.User.ID

	// 2. Normal user tries to CREATE category -> MUST FAIL with 403 Forbidden!
	catPayload := map[string]string{
		"name": fmt.Sprintf("Forbidden Cat %d", ts%10000),
	}
	catBody, _ := json.Marshal(catPayload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewReader(catBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+normalToken)
	createResp, err := testReq(app, createReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, createResp.StatusCode)

	// 3. Generate token for Admin user (is_admin = true)
	adminToken, err := utils.GenerateToken(userID, normalUserEmail, true)
	require.NoError(t, err)

	// 4. Admin creates category -> MUST SUCCEED with 201 Created
	adminCreateReq := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewReader(catBody))
	adminCreateReq.Header.Set("Content-Type", "application/json")
	adminCreateReq.Header.Set("Authorization", "Bearer "+adminToken)
	adminCreateResp, err := testReq(app, adminCreateReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, adminCreateResp.StatusCode)

	var adminCreateRes struct {
		Data model.CategoryResponse `json:"data"`
	}
	json.NewDecoder(adminCreateResp.Body).Decode(&adminCreateRes)
	catID := adminCreateRes.Data.ID

	// 5. Public user (unauthenticated) can GET all categories
	publicListReq := httptest.NewRequest(http.MethodGet, "/api/v1/categories?limit=10&offset=0", nil)
	publicListResp, err := testReq(app, publicListReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, publicListResp.StatusCode)

	// 6. Public user can GET category by ID
	publicGetReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/categories/%d", catID), nil)
	publicGetResp, err := testReq(app, publicGetReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, publicGetResp.StatusCode)

	// 7. Normal user tries to DELETE category -> MUST FAIL with 403 Forbidden!
	delReqNormal := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/categories/%d", catID), nil)
	delReqNormal.Header.Set("Authorization", "Bearer "+normalToken)
	delRespNormal, err := testReq(app, delReqNormal)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, delRespNormal.StatusCode)

	// 8. Admin user deletes category -> MUST SUCCEED with 200 OK
	delReqAdmin := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/categories/%d", catID), nil)
	delReqAdmin.Header.Set("Authorization", "Bearer "+adminToken)
	delRespAdmin, err := testReq(app, delReqAdmin)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, delRespAdmin.StatusCode)

	var delRes helper.BaseResponse
	json.NewDecoder(delRespAdmin.Body).Decode(&delRes)
	assert.Equal(t, 200, delRes.Code)
}
