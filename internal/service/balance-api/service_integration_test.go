package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/models"
	"github.com/yanmoyy/tbi/internal/test"
)

func getTestDBClient(t *testing.T) *database.Client {
	t.Helper()
	test.CheckDBFlag(t)
	cfg := config.LoadWithPath("../../../.env")
	cfg.DB.Host = "localhost"
	c := database.NewClient(cfg.DB)
	return c
}

func getTestService(t *testing.T) *Service {
	t.Helper()
	// Initialize Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := getTestDBClient(t)
	err := db.ClearAll()
	require.NoError(t, err)
	svc := &Service{
		db:     db,
		router: router,
		port:   "8081", // Use different port for tests
	}
	svc.setupRoutes()
	return svc
}

func setupDatasets(t *testing.T, s *Service) {
	balances := []models.TokenBalance{
		{
			Address:   "address1",
			TokenPath: "token1",
			Amount:    100,
		},
		{
			Address:   "address1",
			TokenPath: "token2",
			Amount:    100,
		},
		{
			Address:   "address2",
			TokenPath: "token1",
			Amount:    200,
		},
		{
			Address:   "address3",
			TokenPath: "token2",
			Amount:    300,
		},
	}
	transfers := []models.TokenTransfer{
		{
			FromAddress: "address1",
			ToAddress:   "address2",
			TokenPath:   "token1",
			Amount:      100,
		},
		{
			FromAddress: "address1",
			ToAddress:   "address3",
			TokenPath:   "token2",
			Amount:      200,
		},
	}
	require.NoError(t, s.db.CreateTokenBalanceList(t.Context(), balances))
	require.NoError(t, s.db.CreateTokenTransferList(t.Context(), transfers))
}

func TestAPIService(t *testing.T) {
	s := getTestService(t)
	clear := func() {
		err := s.db.ClearAll()
		require.NoError(t, err)
	}
	clear()
	defer clear()

	setupDatasets(t, s)

	getRecorder := func(endpoint string) *httptest.ResponseRecorder {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", endpoint, nil)
		s.router.ServeHTTP(r, req)
		return r
	}

	t.Run("HealthCheck", func(t *testing.T) {
		r := getRecorder("/health")
		require.Equal(t, http.StatusOK, r.Code)

		var response map[string]string
		require.NoError(t, json.Unmarshal(r.Body.Bytes(), &response))

		require.Equal(t, "ok", response["status"])
	})

	t.Run("GetBalances", func(t *testing.T) {
		r := getRecorder("/tokens/balances")
		require.Equal(t, http.StatusOK, r.Code)

		var resp getBalancesResp
		require.NoError(t, json.Unmarshal(r.Body.Bytes(), &resp))

		require.Len(t, resp.Balances, 2)
		assert.Equal(t, "token1", resp.Balances[0].TokenPath)
		assert.Equal(t, int64(300), resp.Balances[0].Amount)
		assert.Equal(t, "token2", resp.Balances[1].TokenPath)
		assert.Equal(t, int64(400), resp.Balances[1].Amount)
	})

	t.Run("GetBalancesWithAddress", func(t *testing.T) {
		r := getRecorder("/tokens/balances?address=address1")
		require.Equal(t, http.StatusOK, r.Code)

		var resp getBalancesResp
		require.NoError(t, json.Unmarshal(r.Body.Bytes(), &resp))

		require.Len(t, resp.Balances, 2)
		assert.Equal(t, "token1", resp.Balances[0].TokenPath)
		assert.Equal(t, int64(100), resp.Balances[0].Amount)
		assert.Equal(t, "token2", resp.Balances[1].TokenPath)
		assert.Equal(t, int64(100), resp.Balances[1].Amount)
	})

	t.Run("GetTokenAccountBalances", func(t *testing.T) {
		r := getRecorder("/tokens/token1/balances")
		require.Equal(t, http.StatusOK, r.Code)

		var resp getAccountBalancesResp
		require.NoError(t, json.Unmarshal(r.Body.Bytes(), &resp))

		require.Len(t, resp.AccountBalances, 2)
		assert.Equal(t, "token1", resp.AccountBalances[0].TokenPath)
		assert.Equal(t, "token1", resp.AccountBalances[1].TokenPath)
		assert.Equal(t, "address1", resp.AccountBalances[0].Address)
		assert.Equal(t, "address2", resp.AccountBalances[1].Address)
		assert.Equal(t, int64(100), resp.AccountBalances[0].Amount)
		assert.Equal(t, int64(200), resp.AccountBalances[1].Amount)
	})

	t.Run("GetTokenAccountBalancesWithAddress", func(t *testing.T) {
		r := getRecorder("/tokens/token1/balances?address=address1")
		require.Equal(t, http.StatusOK, r.Code)

		var resp getAccountBalancesResp
		require.NoError(t, json.Unmarshal(r.Body.Bytes(), &resp))

		require.Len(t, resp.AccountBalances, 1)
		assert.Equal(t, "token1", resp.AccountBalances[0].TokenPath)
		assert.Equal(t, "address1", resp.AccountBalances[0].Address)
		assert.Equal(t, int64(100), resp.AccountBalances[0].Amount)
	})

	t.Run("GetTransferHistory", func(t *testing.T) {
		r := getRecorder("/tokens/transfer-history")
		require.Equal(t, http.StatusOK, r.Code)

		var resp getTransferHistoryResp
		require.NoError(t, json.Unmarshal(r.Body.Bytes(), &resp))

		require.Len(t, resp.Transfers, 2)
		assert.Equal(t, "address1", resp.Transfers[0].FromAddress)
		assert.Equal(t, "address2", resp.Transfers[0].ToAddress)
		assert.Equal(t, "token1", resp.Transfers[0].TokenPath)
		assert.Equal(t, int64(100), resp.Transfers[0].Amount)
		assert.Equal(t, "address1", resp.Transfers[1].FromAddress)
		assert.Equal(t, "address3", resp.Transfers[1].ToAddress)
		assert.Equal(t, "token2", resp.Transfers[1].TokenPath)
		assert.Equal(t, int64(200), resp.Transfers[1].Amount)
	})

	t.Run("GetTransferHistoryWithAddress", func(t *testing.T) {
		r := getRecorder("/tokens/transfer-history?address=address2")
		require.Equal(t, http.StatusOK, r.Code)

		var resp getTransferHistoryResp
		require.NoError(t, json.Unmarshal(r.Body.Bytes(), &resp))

		require.Len(t, resp.Transfers, 1)
		assert.Equal(t, "address1", resp.Transfers[0].FromAddress)
		assert.Equal(t, "address2", resp.Transfers[0].ToAddress)
		assert.Equal(t, "token1", resp.Transfers[0].TokenPath)
		assert.Equal(t, int64(100), resp.Transfers[0].Amount)
	})

	t.Run("ResponseNotFound", func(t *testing.T) {
		r := getRecorder("/tokens/transfer-history?address=address4")
		require.Equal(t, http.StatusNotFound, r.Code)

		var resp map[string]string
		require.NoError(t, json.Unmarshal(r.Body.Bytes(), &resp))

		require.Equal(t, "no transfer history found", resp["error"])
	})
}
