package indexer

import (
	"context"
	"fmt"

	"github.com/yanmoyy/tbi/internal/models"
)

type getTransactionsResp struct {
	Transactions []Transaction `graphql:"getTransactions"`
}

type GetTransactionsVars struct {
	StartHeight int
	EndHeight   int
	StartIndex  int
	EndIndex    int
}

func (c *Client) GetTransactions(ctx context.Context, vars GetTransactionsVars) (getTransactionsResp, error) {
	var resp getTransactionsResp

	variables := map[string]any{
		"height_eq": vars.StartHeight,
		"height_gt": vars.StartHeight,
		"height_lt": vars.EndHeight,
		"index_eq":  vars.StartIndex,
		"index_gt":  vars.StartIndex,
		"index_lt":  vars.EndIndex,
	}

	err := c.queryTransactions(ctx, &resp, variables)
	if err != nil {
		return getTransactionsResp{}, err
	}
	return resp, err
}

func (c *Client) queryTransactions(ctx context.Context, resp *getTransactionsResp, variables map[string]any) error {
	for _, url := range c.indexerURLs {
		client, ok := c.clients[url]
		if !ok {
			return fmt.Errorf("client not found: %s", url)
		}
		err := client.Exec(ctx, getTransactionsGQL, resp, variables)
		if err != nil {
			return fmt.Errorf("client.Query: %w", err)
		}
		if len(resp.Transactions) > 0 {
			return nil
		}
	}
	return ErrFailedAllEndpoints
}

func (q *getTransactionsResp) ToModel() ([]models.Transaction, error) {
	result := make([]models.Transaction, len(q.Transactions))
	for i, t := range q.Transactions {
		model, err := t.ToModel()
		if err != nil {
			return nil, err
		}
		result[i] = model
	}
	return result, nil
}
