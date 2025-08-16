package indexer

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/hasura/go-graphql-client/pkg/jsonutil"
)

const timeoutPerURL = time.Second * 5

type subscribeBlockQuery struct {
	GetBlocks Block `graphql:" getBlocks(where: {})"`
}

// subscribe to new blocks
func (c *Client) SubscribeBlocks(ctx context.Context, blockChan chan<- Block, done chan<- struct{}) {
	go func() {
		for _, url := range c.indexerURLs {
			wsURL := toWsURL(url)
			slog.Info("subscribing to blocks", "url", wsURL)
			err := c.subscribeToURL(ctx, wsURL, blockChan)
			if err != nil {
				slog.Error("subscribeToURL", "err", err)
			}
		}
		close(done)
	}()
}

// subscribeToURL handles subscription to blocks for a single URL.
func (c *Client) subscribeToURL(ctx context.Context, url string, blockCh chan<- Block) error {
	urlCtx, cancel := context.WithTimeout(context.Background(), timeoutPerURL)
	defer cancel()

	client := graphql.NewSubscriptionClient(url)
	defer client.Close()

	var q subscribeBlockQuery
	errCh := make(chan error)
	dataCh := make(chan []byte)

	_, err := client.Subscribe(&q, nil, func(dataValue []byte, errValue error) error {
		if errValue != nil {
			errCh <- errValue
			return nil
		}
		dataCh <- dataValue
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	go func() {
		if err := client.Run(); err != nil {
			errCh <- fmt.Errorf("failed to run: %w", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-urlCtx.Done():
			return fmt.Errorf("timeout: %w", urlCtx.Err())
		case err := <-errCh:
			return fmt.Errorf("from errorCh: %w", err)
		case data := <-dataCh:
			err := jsonutil.UnmarshalGraphQL(data, &q)
			if err != nil {
				return fmt.Errorf("failed to unmarshal: %w", err)
			}
			urlCtx, cancel = context.WithTimeout(context.Background(), timeoutPerURL)
			defer cancel()
			blockCh <- q.GetBlocks
		}
	}
}

func toWsURL(httpURL string) string {
	return strings.Replace(httpURL, "http", "ws", 1)
}
