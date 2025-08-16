package indexer

import (
	"errors"
	"net/http"
	"strings"

	"github.com/hasura/go-graphql-client"
	"github.com/yanmoyy/tbi/internal/config"
)

var (
	ErrFailedAllEndpoints = errors.New("failed to query from all endpoints")
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
	queryURLs := make([]string, l)
	for i, url := range cfg.IndexerURLs {
		url = getQueryURL(url)
		clients[url] = graphql.NewClient(url, http.DefaultClient)
		queryURLs[i] = url
	}
	return &Client{
		indexerURLs: queryURLs,
		clients:     clients,
	}
}

func getQueryURL(url string) string {
	if strings.HasSuffix(url, "/query") {
		return url
	}
	return url + "/query"
}
