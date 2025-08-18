package processor

import (
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/models"
	"github.com/yanmoyy/tbi/internal/sqs"
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

func TestProcessMessages(t *testing.T) {
	db := getTestDBClient(t)
	s := Service{db: db}

	clear := func() {
		err := s.db.ClearTokenBalances()
		require.NoError(t, err)
		err = s.db.ClearTransactions()
		require.NoError(t, err)
	}

	t.Run("Check Balance", func(t *testing.T) {
		clear()
		defer clear()

		type input struct {
			events []models.TransferEvent
		}
		type expected struct {
			balances []models.TokenBalance
		}

		ctx := t.Context()
		tester := func(t *testing.T, i input, e expected) {
			msgs := make([]sqs.Message, len(i.events))
			for i, evt := range i.events {
				body, err := json.Marshal(evt)
				require.NoError(t, err)
				msgs[i] = sqs.Message{Body: string(body)}
			}
			err := s.processMessages(ctx, msgs)
			require.NoError(t, err)

			balances, err := s.db.GetTokenBalanceList(ctx)
			require.NoError(t, err)

			slog.Info("Balances", "balances", balances)

			require.Len(t, balances, len(e.balances))

			for i, b := range balances {
				require.Equal(t, e.balances[i].Address, b.Address)
				require.Equal(t, e.balances[i].TokenPath, b.TokenPath)
				require.Equal(t, e.balances[i].Amount, b.Amount)
			}
		}
		// Test Transfer
		tester(t, input{
			events: []models.TransferEvent{
				{
					Func:    models.Mint,
					PkgPath: "token1",
					To:      "a",
					Value:   20,
				},
				{
					Func:    models.Transfer,
					PkgPath: "token1",
					From:    "a",
					To:      "b",
					Value:   5,
				},
				{
					Func:    models.Burn,
					PkgPath: "token1",
					From:    "a",
					Value:   5,
				},
				{
					Func:    models.Transfer,
					PkgPath: "token1",
					From:    "b",
					To:      "a",
					Value:   2,
				},
			},
		}, expected{
			balances: []models.TokenBalance{
				{
					Address:   "a",
					TokenPath: "token1",
					Amount:    12,
				},
				{
					Address:   "b",
					TokenPath: "token1",
					Amount:    3,
				},
			},
		})
	})
}
