package indexer

import (
	"context"
	"fmt"

	"github.com/yanmoyy/tbi/internal/models"
)

type getBlocksResp struct {
	Blocks []Block `graphql:"getBlocks"`
}

func (c *Client) GetBlocks(ctx context.Context, heightStart, heightEnd int) (getBlocksResp, error) {

	var resp getBlocksResp
	vars := map[string]any{
		"height_eq": heightStart,
		"height_gt": heightStart,
		"height_lt": heightEnd,
	}
	err := c.queryBlocks(ctx, &resp, vars)
	if err != nil {
		return getBlocksResp{}, err
	}
	return resp, nil
}

func (c *Client) queryBlocks(ctx context.Context, resp *getBlocksResp, variables map[string]any) error {
	for _, url := range c.indexerURLs {
		client, ok := c.clients[url]
		if !ok {
			return fmt.Errorf("client not found: %s", url)
		}
		if err := client.WithDebug(true).Exec(ctx, getBlocksGQL, resp, variables); err != nil {
			return fmt.Errorf("querying %s: %w", url, err)
		}
		if len(resp.Blocks) > 0 {
			return nil
		}
	}
	return ErrFailedAllEndpoints
}

func (r *getBlocksResp) ToModel() []models.Block {
	result := make([]models.Block, len(r.Blocks))
	for i, b := range r.Blocks {
		result[i] = b.ToModel()
	}
	return result
}
