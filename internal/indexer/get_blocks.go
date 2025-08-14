package indexer

import (
	"context"
	"log/slog"

	"github.com/yanmoyy/tbi/internal/models"
)

type getBlockQuery struct {
	GetBlocks []block `graphql:"getBlocks(where: { height: { gt: $height_gt, lt: $height_lt } })"`
}

type GetBlocksFilter struct {
	HeightGT int
	HeightLT int
}

func (c *Client) GetBlocks(ctx context.Context, filter GetBlocksFilter) ([]models.Block, error) {

	var q getBlockQuery

	variables := map[string]any{
		"height_gt": filter.HeightGT,
		"height_lt": filter.HeightLT,
	}

	for _, url := range c.indexerURLs {
		client, ok := c.clients[url]
		if !ok {
			continue
		}
		err := client.Query(ctx, &q, variables)
		if err != nil {
			slog.Error("failed to query blocks", "err", err, "url", url)
			continue
		}
		if len(q.GetBlocks) == 0 {
			slog.Error("no blocks found", "url", url)
			continue
		}
		return convert(q.GetBlocks, blockConvertor)
	}

	return nil, ErrFailedAllEndpoints
}
