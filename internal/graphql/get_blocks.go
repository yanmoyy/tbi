package graphql

import (
	"context"
	"log/slog"

	"github.com/yanmoyy/tbi/internal/models"
)

type getBlockQuery struct {
	GetBlocks []block `graphql:"getBlocks(where: { height: { gt: $height_gt, lt: $height_lt } })"`
}

type GetBlocksFilter struct {
	HeightGT int64
	HeightLT int64
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
		return convert(q.GetBlocks, blockToModel), nil
	}

	return nil, ErrFailedAllEndpoints
}

func blockToModel(b block) models.Block { return b.toModel() }

func (b *block) toModel() models.Block {
	return models.Block{
		Hash:     b.Hash,
		Height:   b.Height,
		Time:     b.Time,
		NumTxs:   b.NumTxs,
		TotalTxs: b.TotalTxs,
	}
}
