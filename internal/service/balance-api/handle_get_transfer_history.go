package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yanmoyy/tbi/internal/models"
)

type transfer struct {
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	TokenPath   string `json:"tokenPath"`
	Amount      int64  `json:"amount"`
}

type getTransferHistoryResp struct {
	Transfers []transfer `json:"transfers"`
}

func (s *Service) handleGetTransferHistory(c *gin.Context) {
	ctx := c.Request.Context()
	address := c.Query("address")

	var resp []models.TokenTransfer
	var err error
	if address == "" {
		resp, err = s.db.GetTokenTransferList(ctx)
	} else {
		resp, err = s.db.GetTokenTransferListWithAddress(ctx, address)
	}
	if err != nil {
		respondServerError(c, "failed to query transfer history", err)
		return
	}
	slog.Info("handleGetTransferHistory", "resp", resp)
	if len(resp) == 0 {
		respondNotFound(c, "no transfer history found")
		return
	}

	transfers := make([]transfer, len(resp))
	for i, t := range resp {
		transfers[i] = transfer{
			FromAddress: t.FromAddress,
			ToAddress:   t.ToAddress,
			TokenPath:   t.TokenPath,
			Amount:      t.Amount,
		}
	}
	c.JSON(http.StatusOK, getTransferHistoryResp{Transfers: transfers})
}
