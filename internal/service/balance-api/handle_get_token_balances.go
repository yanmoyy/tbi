package api

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/yanmoyy/tbi/internal/models"
)

type tokenBalance struct {
	TokenPath string `json:"tokenPath"`
	Amount    int64  `json:"amount"`
}

type getBalancesResp struct {
	Balances []tokenBalance `json:"balances"`
}

func (s *Service) handleGetTokenBalances(c *gin.Context) {
	ctx := c.Request.Context()
	address := c.Query("address")
	var resp []models.TokenBalance
	var err error
	if address == "" {
		resp, err = s.db.GetTokenBalanceList(ctx)
	} else {
		resp, err = s.db.GetTokenBalanceListWithAddress(ctx, address)
	}
	if err != nil {
		respondServerError(c, "failed to query balances", err)
		return
	}
	if len(resp) == 0 {
		respondNotFound(c, "no balances found")
		return
	}
	// sum up balances with same tokenPath
	sum := make(map[string]int64)
	for _, b := range resp {
		sum[b.TokenPath] += b.Amount
	}
	balances := make([]tokenBalance, 0, len(sum))
	for k, v := range sum {
		balances = append(balances, tokenBalance{
			TokenPath: k,
			Amount:    v,
		})
	}
	// sort by tokenPath
	sort.Slice(balances, func(i, j int) bool {
		return balances[i].TokenPath < balances[j].TokenPath
	})
	c.JSON(http.StatusOK, getBalancesResp{Balances: balances})
}
