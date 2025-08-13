package graphql

import (
	"context"
	"flag"
	"testing"

	"github.com/yanmoyy/tbi/internal/config"
)

func TestGetTransactions(t *testing.T) {
	flag.Parse()
	if *offline {
		t.Skip("Skipping test in offline mode")
	}
	cfg := config.GraphQL{
		IndexerURLs: []string{indexerURL},
	}
	c := NewClient(cfg)
	transactions, err := c.GetTransactions(context.Background(), GetTransactionsFilter{
		BlockHeightGT: 0,
		BlockHeightLT: 2000,
		IndexLT:       1,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("total transactions: %d", len(transactions))
}
