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
	"evermos-api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionHandler_E2E(t *testing.T) {
	app := setupTestApp(t)

	ts := time.Now().UnixNano()

	// 1. Register Seller
	sellerEmail := fmt.Sprintf("seller_tx_%d@evermos.com", ts)
	sellerPhone := fmt.Sprintf("0874%d", ts%100000000)
	regSeller := map[string]string{
		"name":     "Seller TX",
		"email":    sellerEmail,
		"phone":    sellerPhone,
		"password": "Password123!",
	}
	bodySeller, _ := json.Marshal(regSeller)
	reqRegSeller := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(bodySeller))
	reqRegSeller.Header.Set("Content-Type", "application/json")
	respRegSeller, err := testReq(app, reqRegSeller)
	require.NoError(t, err)

	var resSeller struct {
		Data model.AuthResponse `json:"data"`
	}
	json.NewDecoder(respRegSeller.Body).Decode(&resSeller)
	sellerToken := resSeller.Data.Token
	sellerID := resSeller.Data.User.ID

	// 2. Admin creates Category
	adminToken, err := utils.GenerateToken(sellerID, sellerEmail, true)
	require.NoError(t, err)

	catPayload := map[string]string{
		"name": fmt.Sprintf("Fashion_%d", ts%10000),
	}
	catBody, _ := json.Marshal(catPayload)
	reqCat := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewReader(catBody))
	reqCat.Header.Set("Content-Type", "application/json")
	reqCat.Header.Set("Authorization", "Bearer "+adminToken)
	respCat, err := testReq(app, reqCat)
	require.NoError(t, err)

	var resCat struct {
		Data model.CategoryResponse `json:"data"`
	}
	json.NewDecoder(respCat.Body).Decode(&resCat)
	categoryID := resCat.Data.ID

	// 3. Seller creates Product
	prodPayload := map[string]interface{}{
		"category_id": categoryID,
		"name":        "Evermos Batik Shirt",
		"description": "Original Indonesian Batik",
		"price":       150000,
		"quantity":    8,
	}
	bodyProd, _ := json.Marshal(prodPayload)
	reqProd := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(bodyProd))
	reqProd.Header.Set("Content-Type", "application/json")
	reqProd.Header.Set("Authorization", "Bearer "+sellerToken)
	respProd, err := testReq(app, reqProd)
	require.NoError(t, err)

	var resProd struct {
		Data model.ProductResponse `json:"data"`
	}
	json.NewDecoder(respProd.Body).Decode(&resProd)
	productID := resProd.Data.ID

	// 4. Register Buyer
	buyerEmail := fmt.Sprintf("buyer_tx_%d@evermos.com", ts)
	buyerPhone := fmt.Sprintf("0875%d", ts%100000000)
	regBuyer := map[string]string{
		"name":     "Buyer TX",
		"email":    buyerEmail,
		"phone":    buyerPhone,
		"password": "Password123!",
	}
	bodyBuyer, _ := json.Marshal(regBuyer)
	reqRegBuyer := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(bodyBuyer))
	reqRegBuyer.Header.Set("Content-Type", "application/json")
	respRegBuyer, err := testReq(app, reqRegBuyer)
	require.NoError(t, err)

	var resBuyer struct {
		Data model.AuthResponse `json:"data"`
	}
	json.NewDecoder(respRegBuyer.Body).Decode(&resBuyer)
	buyerToken := resBuyer.Data.Token

	// 5. Buyer creates Address
	addrPayload := map[string]interface{}{
		"judul_alamat":   "Alamat Pengiriman",
		"penerima_nama":  "Buyer TX",
		"penerima_phone": buyerPhone,
		"provinsi":       "DKI Jakarta",
		"provinsi_id":    "31",
		"kabupaten":      "Jakarta Timur",
		"kabupaten_id":   "3175",
		"kecamatan":      "Jatinegara",
		"kecamatan_id":   "317503",
		"kelurahan":      "Bidara Cina",
		"kelurahan_id":   "3175031001",
		"detail_alamat":  "Jl. Otista No. 42",
		"is_default":     true,
	}
	bodyAddr, _ := json.Marshal(addrPayload)
	reqAddr := httptest.NewRequest(http.MethodPost, "/api/v1/addresses", bytes.NewReader(bodyAddr))
	reqAddr.Header.Set("Content-Type", "application/json")
	reqAddr.Header.Set("Authorization", "Bearer "+buyerToken)
	respAddr, err := testReq(app, reqAddr)
	require.NoError(t, err)

	var resAddr struct {
		Data model.AddressResponse `json:"data"`
	}
	json.NewDecoder(respAddr.Body).Decode(&resAddr)
	addressID := resAddr.Data.ID

	// 6. Buyer executes Checkout (POST /api/v1/transactions)
	txPayload := map[string]interface{}{
		"address_id": addressID,
		"products": []map[string]interface{}{
			{
				"product_id": productID,
				"quantity":   2,
			},
		},
	}
	bodyTx, _ := json.Marshal(txPayload)
	reqTx := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewReader(bodyTx))
	reqTx.Header.Set("Content-Type", "application/json")
	reqTx.Header.Set("Authorization", "Bearer "+buyerToken)
	respTx, err := testReq(app, reqTx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, respTx.StatusCode)

	var resTx struct {
		Data model.TransactionResponse `json:"data"`
	}
	json.NewDecoder(respTx.Body).Decode(&resTx)
	transactionID := resTx.Data.ID
	assert.Equal(t, float64(300000), resTx.Data.TotalPrice)
	assert.NotEmpty(t, resTx.Data.InvoiceNumber)
	assert.Len(t, resTx.Data.ProductLogs, 1)

	// 7. Buyer gets personal transactions (GET /api/v1/transactions/me)
	reqMyTx := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/me", nil)
	reqMyTx.Header.Set("Authorization", "Bearer "+buyerToken)
	respMyTx, err := testReq(app, reqMyTx)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respMyTx.StatusCode)

	// 8. Seller (another user) tries to view Buyer's transaction -> MUST RETURN 404
	reqTamper := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/transactions/me/%d", transactionID), nil)
	reqTamper.Header.Set("Authorization", "Bearer "+sellerToken)
	respTamper, err := testReq(app, reqTamper)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, respTamper.StatusCode)

	// 9. Buyer updates status to "completed"
	statusPayload := map[string]string{"status": "completed"}
	bodyStatus, _ := json.Marshal(statusPayload)
	reqStatus := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/transactions/me/%d", transactionID), bytes.NewReader(bodyStatus))
	reqStatus.Header.Set("Content-Type", "application/json")
	reqStatus.Header.Set("Authorization", "Bearer "+buyerToken)
	respStatus, err := testReq(app, reqStatus)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respStatus.StatusCode)

	// 10. Admin views all transactions (GET /api/v1/transactions)
	reqAdminList := httptest.NewRequest(http.MethodGet, "/api/v1/transactions", nil)
	reqAdminList.Header.Set("Authorization", "Bearer "+adminToken)
	respAdminList, err := testReq(app, reqAdminList)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respAdminList.StatusCode)
}
