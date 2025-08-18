package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanmoyy/tbi/internal/models"
)

func TestCreateTranscations(t *testing.T) {
	c := getTestClient(t)
	blocks := []models.Block{
		{
			Hash:   "hash1",
			Height: 0,
		},
		{
			Hash:   "hash2",
			Height: 1,
		},
	}
	transactions := []models.Transaction{
		{
			Hash:        "hash1",
			BlockHeight: 0,
			GasFee:      []byte("{}"),
			Messages:    []byte("[]"),
			Response:    []byte("{}"),
		},
		{
			Hash:        "hash2",
			BlockHeight: 1,
			GasFee:      []byte("{}"),
			Messages:    []byte("[]"),
			Response:    []byte("{}"),
		},
	}

	ctx := context.Background()

	err := c.ClearTransactions()
	require.NoError(t, err)

	err = c.CreateBlockList(ctx, blocks)
	require.NoError(t, err)
	err = c.CreateTransactionList(ctx, transactions)
	require.NoError(t, err)

	err = c.ClearTransactions()
	require.NoError(t, err)
}
