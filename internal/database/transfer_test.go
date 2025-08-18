package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanmoyy/tbi/internal/models"
)

func TestCreateTransfer(t *testing.T) {
	c := getTestClient(t)
	clear := func() {
		err := c.ClearTokenTransfers()
		require.NoError(t, err)
	}
	clear()
	defer clear()

	ctx := context.Background()

	err := c.CreateTokenTransfer(ctx, models.TokenTransfer{
		FromAddress: "from",
		ToAddress:   "to",
		Amount:      100,
		TokenPath:   "token",
	})
	require.NoError(t, err)

	transfers, err := c.GetTokenTransferList(ctx)
	require.NoError(t, err)
	require.Len(t, transfers, 1)
	require.Equal(t, "from", transfers[0].FromAddress)
	require.Equal(t, "to", transfers[0].ToAddress)
	require.Equal(t, int64(100), transfers[0].Amount)
	require.Equal(t, "token", transfers[0].TokenPath)
	require.Equal(t, uint(1), transfers[0].ID)
}
