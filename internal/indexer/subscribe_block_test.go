package indexer

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/models"
	"github.com/yanmoyy/tbi/internal/test"
)

var minBlocks = flag.Int("min-blocks", 0, "Minimum number of blocks to receive before exiting")

func TestSubscribeBlocks(t *testing.T) {
	test.CheckIndexerFlag(t)
	if *minBlocks == 0 {
		t.Skip("Skipping test without minimum blocks")
	}
	cfg := config.GraphQL{
		IndexerURLs: []string{indexerURL},
	}
	c := NewClient(cfg)
	blockChan := make(chan models.Block)
	done := make(chan struct{})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	blockCnt := 0
	defer func() {
		if blockCnt < *minBlocks {
			t.Errorf("Expected at least %d blocks, received %d", *minBlocks, blockCnt)
		}
	}()

	// Start SubscribeBlocks
	c.SubscribeBlocks(ctx, blockChan, done)

	for {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		case block := <-blockChan:
			blockCnt++
			t.Logf("Block #%d: %+v", blockCnt, block)
			if blockCnt >= *minBlocks {
				return
			}
		}
	}
}
