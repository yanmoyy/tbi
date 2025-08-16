package indexer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/test"
)

func TestGetTransactions(t *testing.T) {
	test.CheckIndexerFlag(t)
	cfg := config.GraphQL{
		IndexerURLs: []string{indexerURL},
	}
	c := NewClient(cfg)
	start := TxIndex{BlockHeight: 0, Index: 0}
	end := TxIndex{BlockHeight: 1, Index: 1}
	resp, err := c.GetTransactions(context.Background(), start, end)
	require.NoError(t, err)
	t.Logf("total transactions: %d", len(resp.Transactions))
}
