package graphql

import (
	"context"
	"testing"
	"time"

	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/models"
)

func TestSubscribeBlocks(t *testing.T) {
	cfg := config.GraphQL{
		IndexerURLs: []string{indexerURL, "https://indexer.onbloc.xyz/graphql/query"},
	}
	c := NewClient(cfg)
	blockChan := make(chan models.Block)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	blockCnt := 0
	completeCh := make(chan struct{})
	c.SubscribeBlocks(ctx, blockChan, completeCh)
	for {
		select {
		case <-ctx.Done():
			t.Error("timeout")
			return
		case <-completeCh:
			if blockCnt == 0 {
				t.Error("no blocks received")
			}
			return
		case block := <-blockChan:
			blockCnt += 1
			t.Logf("block #%d: %+v", blockCnt, block)
		}
	}
}
