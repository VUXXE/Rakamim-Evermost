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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddressHandler_E2E(t *testing.T) {
	app := setupTestApp(t)

	ts := time.Now().UnixNano()
	userAEmail := fmt.Sprintf("addr_a_%d@evermos.com", ts)
	userAPhone := fmt.Sprintf("0894%d", ts%100000000)

	// 1. Register User A
	regPayloadA := map[string]string{
		"name":     "Address Tester A",
		"email":    userAEmail,
		"phone":    userAPhone,
		"password": "Password123!",
	}
	bodyA, _ := json.Marshal(regPayloadA)
	regReqA := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(bodyA))
	regReqA.Header.Set("Content-Type", "application/json")
	regRespA, err := testReq(app, regReqA)
	require.NoError(t, err)

	var resA struct {
		Data model.AuthResponse `json:"data"`
	}
	json.NewDecoder(regRespA.Body).Decode(&resA)
	tokenA := resA.Data.Token

	// 2. Register User B
	userBEmail := fmt.Sprintf("addr_b_%d@evermos.com", ts)
	userBPhone := fmt.Sprintf("0895%d", ts%100000000)
	regPayloadB := map[string]string{
		"name":     "Address Tester B",
		"email":    userBEmail,
		"phone":    userBPhone,
		"password": "Password123!",
	}
	bodyB, _ := json.Marshal(regPayloadB)
	regReqB := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(bodyB))
	regReqB.Header.Set("Content-Type", "application/json")
	regRespB, err := testReq(app, regReqB)
	require.NoError(t, err)

	var resB struct {
		Data model.AuthResponse `json:"data"`
	}
	json.NewDecoder(regRespB.Body).Decode(&resB)
	tokenB := resB.Data.Token

	// 3. User A creates an address
	addrPayload := map[string]interface{}{
		"judul_alamat":   "Rumah Bandung",
		"penerima_nama":  "Penerima A",
		"penerima_phone": "0894000000",
		"provinsi":       "Jawa Barat",
		"provinsi_id":    "32",
		"kabupaten":      "Bandung",
		"kabupaten_id":   "3273",
		"kecamatan":      "Coblong",
		"kecamatan_id":   "327301",
		"kelurahan":      "Dago",
		"kelurahan_id":   "3273011001",
		"detail_alamat":  "Jl. Ir. H. Juanda No. 100",
		"is_default":     true,
	}
	bodyAddr, _ := json.Marshal(addrPayload)
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/v1/addresses", bytes.NewReader(bodyAddr))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+tokenA)

	respCreate, err := testReq(app, reqCreate)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, respCreate.StatusCode)

	var createRes struct {
		Data model.AddressResponse `json:"data"`
	}
	json.NewDecoder(respCreate.Body).Decode(&createRes)
	addrID := createRes.Data.ID

	// 4. User A gets personal addresses
	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/addresses?limit=10&offset=0", nil)
	reqList.Header.Set("Authorization", "Bearer "+tokenA)
	respList, err := testReq(app, reqList)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respList.StatusCode)

	// 5. User B tries to view User A's address -> MUST return 404 Not Found!
	reqTamperGet := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/addresses/%d", addrID), nil)
	reqTamperGet.Header.Set("Authorization", "Bearer "+tokenB)
	respTamperGet, err := testReq(app, reqTamperGet)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, respTamperGet.StatusCode)

	// 6. User B tries to delete User A's address -> MUST return 404 Not Found!
	reqTamperDel := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/addresses/%d", addrID), nil)
	reqTamperDel.Header.Set("Authorization", "Bearer "+tokenB)
	respTamperDel, err := testReq(app, reqTamperDel)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, respTamperDel.StatusCode)

	// 7. User A legitimately deletes own address -> MUST return 200 OK
	reqDel := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/addresses/%d", addrID), nil)
	reqDel.Header.Set("Authorization", "Bearer "+tokenA)
	respDel, err := testReq(app, reqDel)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respDel.StatusCode)

	var delRes helper.BaseResponse
	json.NewDecoder(respDel.Body).Decode(&delRes)
	assert.Equal(t, 200, delRes.Code)
	assert.Equal(t, "address deleted successfully", delRes.Message)
}
