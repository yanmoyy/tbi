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
	resp, err := c.GetTransactions(context.Background(), GetTransactionsVars{
		StartHeight: 0,
		EndHeight:   1000,
		StartIndex:  0,
		EndIndex:    100,
	})
	require.NoError(t, err)
	t.Logf("total transactions: %d", len(resp.Transactions))
}
