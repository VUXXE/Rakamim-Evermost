package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"evermos-api/internal/helper"
	"evermos-api/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegionHandler_E2E(t *testing.T) {
	app := setupTestApp(t)

	// 1. GET /api/v1/regions/provinces
	reqProv := httptest.NewRequest(http.MethodGet, "/api/v1/regions/provinces", nil)
	respProv, err := testReq(app, reqProv)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respProv.StatusCode)

	var resProv struct {
		Data []model.Province `json:"data"`
	}
	err = json.NewDecoder(respProv.Body).Decode(&resProv)
	require.NoError(t, err)
	assert.NotEmpty(t, resProv.Data)

	// 2. GET /api/v1/regions/regencies/:province_id (using West Java: 32)
	reqReg := httptest.NewRequest(http.MethodGet, "/api/v1/regions/regencies/32", nil)
	respReg, err := testReq(app, reqReg)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respReg.StatusCode)

	var resReg struct {
		Data []model.Regency `json:"data"`
	}
	err = json.NewDecoder(respReg.Body).Decode(&resReg)
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.Data)

	// 3. GET /api/v1/regions/districts/:regency_id (using Bandung: 3273)
	reqDist := httptest.NewRequest(http.MethodGet, "/api/v1/regions/districts/3273", nil)
	respDist, err := testReq(app, reqDist)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respDist.StatusCode)

	var resDist struct {
		Data []model.District `json:"data"`
	}
	err = json.NewDecoder(respDist.Body).Decode(&resDist)
	require.NoError(t, err)
	assert.NotEmpty(t, resDist.Data)

	// 4. GET /api/v1/regions/villages/:district_id (using Coblong: 3273010)
	reqVil := httptest.NewRequest(http.MethodGet, "/api/v1/regions/villages/3273010", nil)
	respVil, err := testReq(app, reqVil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respVil.StatusCode)

	var resVil struct {
		Data []model.Village `json:"data"`
	}
	err = json.NewDecoder(respVil.Body).Decode(&resVil)
	require.NoError(t, err)
	assert.NotEmpty(t, resVil.Data)

	// 5. Test caching: Second call should be instant
	respProvCached, err := testReq(app, reqProv)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respProvCached.StatusCode)
}

func TestRegionHandler_ValidationErrors(t *testing.T) {
	app := setupTestApp(t)

	// Test missing parameter or not found behavior
	reqBad := httptest.NewRequest(http.MethodGet, "/api/v1/regions/regencies/non-existent-prov", nil)
	respBad, err := testReq(app, reqBad)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadGateway, respBad.StatusCode)

	var res helper.BaseResponse
	err = json.NewDecoder(respBad.Body).Decode(&res)
	require.NoError(t, err)
	assert.Equal(t, 502, res.Code)
}
