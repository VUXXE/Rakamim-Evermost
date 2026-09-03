package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
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

func TestProductHandler_MultipartUpload_AndIsolation_E2E(t *testing.T) {
	app := setupTestApp(t)

	ts := time.Now().UnixNano()

	// 1. Register Seller A
	userAEmail := fmt.Sprintf("seller_a_%d@evermos.com", ts)
	userAPhone := fmt.Sprintf("0899%d", ts%100000000)
	regPayloadA := map[string]string{
		"name":     "Seller A",
		"email":    userAEmail,
		"phone":    userAPhone,
		"password": "Password123!",
	}
	bodyA, _ := json.Marshal(regPayloadA)
	regReqA := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(bodyA))
	regReqA.Header.Set("Content-Type", "application/json")
	regRespA, err := testReq(app, regReqA)
	require.NoError(t, err)

	var regResA struct {
		Data model.AuthResponse `json:"data"`
	}
	json.NewDecoder(regRespA.Body).Decode(&regResA)
	tokenA := regResA.Data.Token
	userAID := regResA.Data.User.ID

	// 2. Register Seller B
	userBEmail := fmt.Sprintf("seller_b_%d@evermos.com", ts)
	userBPhone := fmt.Sprintf("0890%d", ts%100000000)
	regPayloadB := map[string]string{
		"name":     "Seller B",
		"email":    userBEmail,
		"phone":    userBPhone,
		"password": "Password123!",
	}
	bodyB, _ := json.Marshal(regPayloadB)
	regReqB := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(bodyB))
	regReqB.Header.Set("Content-Type", "application/json")
	regRespB, err := testReq(app, regReqB)
	require.NoError(t, err)

	var regResB struct {
		Data model.AuthResponse `json:"data"`
	}
	json.NewDecoder(regRespB.Body).Decode(&regResB)
	tokenB := regResB.Data.Token

	// 3. Admin creates Category
	adminToken, err := utils.GenerateToken(userAID, userAEmail, true)
	require.NoError(t, err)

	catPayload := map[string]string{
		"name": fmt.Sprintf("Gadget_%d", ts%10000),
	}
	catBody, _ := json.Marshal(catPayload)
	createCatReq := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewReader(catBody))
	createCatReq.Header.Set("Content-Type", "application/json")
	createCatReq.Header.Set("Authorization", "Bearer "+adminToken)
	catResp, err := testReq(app, createCatReq)
	require.NoError(t, err)

	var catRes struct {
		Data model.CategoryResponse `json:"data"`
	}
	json.NewDecoder(catResp.Body).Decode(&catRes)
	categoryID := catRes.Data.ID

	// 4. Seller A creates Product with Multipart Image Upload
	bodyBuf := &bytes.Buffer{}
	mpWriter := multipart.NewWriter(bodyBuf)
	_ = mpWriter.WriteField("category_id", fmt.Sprintf("%d", categoryID))
	_ = mpWriter.WriteField("name", "Smartphone Flagship Ultra")
	_ = mpWriter.WriteField("description", "Best camera smartphone")
	_ = mpWriter.WriteField("price", "12000000")
	_ = mpWriter.WriteField("quantity", "15")

	fileWriter, err := mpWriter.CreateFormFile("image", "phone.png")
	require.NoError(t, err)
	_, _ = fileWriter.Write([]byte("fake-png-content"))
	_ = mpWriter.Close()

	prodReq := httptest.NewRequest(http.MethodPost, "/api/v1/products", bodyBuf)
	prodReq.Header.Set("Content-Type", mpWriter.FormDataContentType())
	prodReq.Header.Set("Authorization", "Bearer "+tokenA)

	prodResp, err := testReq(app, prodReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, prodResp.StatusCode)

	var prodRes struct {
		Data model.ProductResponse `json:"data"`
	}
	json.NewDecoder(prodResp.Body).Decode(&prodRes)
	prodID := prodRes.Data.ID
	assert.NotEmpty(t, prodRes.Data.ImageURL)
	assert.Contains(t, prodRes.Data.ImageURL, "phone.png")

	// 5. Seller B tries to update Seller A's product -> MUST FAIL with 403 Forbidden!
	patchBuf := &bytes.Buffer{}
	mpPatch := multipart.NewWriter(patchBuf)
	_ = mpPatch.WriteField("price", "1000")
	_ = mpPatch.Close()

	tamperReq := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/products/%d", prodID), patchBuf)
	tamperReq.Header.Set("Content-Type", mpPatch.FormDataContentType())
	tamperReq.Header.Set("Authorization", "Bearer "+tokenB)
	tamperResp, err := testReq(app, tamperReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, tamperResp.StatusCode)

	// 6. Public user browses products with search & category filter
	searchReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products?search=Smartphone&category_id=%d&sort=price_desc", categoryID), nil)
	searchResp, err := testReq(app, searchReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, searchResp.StatusCode)

	// 7. Seller A views personal products via /api/v1/products/me
	myProdsReq := httptest.NewRequest(http.MethodGet, "/api/v1/products/me", nil)
	myProdsReq.Header.Set("Authorization", "Bearer "+tokenA)
	myProdsResp, err := testReq(app, myProdsReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, myProdsResp.StatusCode)

	// 8. Seller A deletes own product
	delReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/products/%d", prodID), nil)
	delReq.Header.Set("Authorization", "Bearer "+tokenA)
	delResp, err := testReq(app, delReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, delResp.StatusCode)

	var delRes helper.BaseResponse
	json.NewDecoder(delResp.Body).Decode(&delRes)
	assert.Equal(t, 200, delRes.Code)
}
