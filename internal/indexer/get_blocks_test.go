package indexer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/test"
)

func TestGetBlocks(t *testing.T) {
	test.CheckIndexerFlag(t)
	cfg := config.GraphQL{
		IndexerURLs: []string{indexerURL},
	}
	c := NewClient(cfg)
	resp, err := c.GetBlocks(context.Background(), 0, 10)
	require.NoError(t, err)
	require.Equal(t, 10, len(resp.Blocks))

	for _, b := range resp.Blocks {
		if b.Height < 0 || b.Height >= 10 {
			t.Fatal("bad block height", b.Height)
		}
	}
}
