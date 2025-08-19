package api

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/yanmoyy/tbi/internal/models"
)

type accountBalance struct {
	Address   string `json:"address"`
	TokenPath string `json:"tokenPath"`
	Amount    int64  `json:"amount"`
}

type getAccountBalancesResp struct {
	AccountBalances []accountBalance `json:"accountBalances"`
}

func (s *Service) handleGetTokenAccountBalances(c *gin.Context) {
	ctx := c.Request.Context()
	tokenPath := c.Param("tokenPath")
	address := c.Query("address")

	var balances []models.TokenBalance
	if address == "" {
		// get all account balances
		resp, err := s.db.GetTokenBalanceListWithToken(ctx, tokenPath)
		if err != nil {
			respondServerError(c, "failed to get token balances", err)
			return
		}
		balances = resp
	} else {
		// get only one account balance
		resp, err := s.db.GetTokenBalance(ctx, address, tokenPath)
		if err != nil {
			respondServerError(c, "failed to get token balances", err)
			return
		}
		balances = append(balances, resp)
	}

	if len(balances) == 0 {
		respondNotFound(c, "no balances found")
		return
	}

	accountBalances := make([]accountBalance, len(balances))
	for i, b := range balances {
		accountBalances[i] = accountBalance{
			Address:   b.Address,
			TokenPath: b.TokenPath,
			Amount:    b.Amount,
		}
	}
	// sort by account address
	sort.Slice(accountBalances, func(i, j int) bool {
		return accountBalances[i].Address < accountBalances[j].Address
	})
	c.JSON(http.StatusOK, getAccountBalancesResp{AccountBalances: accountBalances})
}
