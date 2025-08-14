package indexer

import (
	"context"
	"testing"

	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/test"
)

func TestGetBlocks(t *testing.T) {
	test.CheckIndexerFlag(t)
	cfg := config.GraphQL{
		IndexerURLs: []string{indexerURL},
	}
	c := NewClient(cfg)
	blocks, err := c.GetBlocks(context.Background(), GetBlocksFilter{
		HeightGT: 0,
		HeightLT: 10,
	})

	for _, b := range blocks {
		if b.Height == 0 || b.Height > 10 {
			t.Fatal("bad block height")
		}
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", blocks)
}
