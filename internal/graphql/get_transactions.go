package graphql

import (
	"context"
	"log/slog"

	"github.com/yanmoyy/tbi/internal/models"
)

type getTransactionsQuery struct {
	GetTransactions []transaction `graphql:"getTransactions(where: { block_height: { gt: $height_gt, lt: $height_lt }, index: { lt: $index_lt } })"`
}

type GetTransactionsFilter struct {
	BlockHeightGT int64
	BlockHeightLT int64
	IndexLT       int64
}

func (c *Client) GetTransactions(ctx context.Context, filter GetTransactionsFilter) ([]models.Transaction, error) {

	var q getTransactionsQuery

	variables := map[string]any{
		"height_gt": filter.BlockHeightGT,
		"height_lt": filter.BlockHeightLT,
		"index_lt":  filter.IndexLT,
	}

	for _, url := range c.indexerURLs {
		client, ok := c.clients[url]
		if !ok {
			continue
		}
		err := client.Query(ctx, &q, variables)
		if err != nil {
			slog.Error("failed to query transactions", "err", err, "url", url)
			continue
		}
		if len(q.GetTransactions) == 0 {
			slog.Error("no transactions found", "url", url)
			continue
		}
		return convert(q.GetTransactions, transactionToModel), nil
	}
	return nil, ErrFailedAllEndpoints
}

func transactionToModel(t transaction) models.Transaction {
	return models.Transaction{
		Hash:        t.Hash,
		Success:     t.Success,
		BlockHeight: t.BlockHeight,
		GasUsed:     t.GasUsed,
	}
}
