package indexer

import (
	"net/http"

	"github.com/hasura/go-graphql-client"
	"github.com/yanmoyy/tbi/internal/config"
)

type Client struct {
	indexerURLs []string
	clients     map[string]*graphql.Client
}

func NewClient(cfg config.GraphQL) *Client {
	if len(cfg.IndexerURLs) == 0 {
		panic("no indexer URLs provided")
	}
	l := len(cfg.IndexerURLs)
	clients := make(map[string]*graphql.Client, l)
	for _, url := range cfg.IndexerURLs {
		clients[url] = graphql.NewClient(url, http.DefaultClient)
	}
	return &Client{
		indexerURLs: cfg.IndexerURLs,
		clients:     clients,
	}
}
