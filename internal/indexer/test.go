package indexer

import "errors"

// Endpoint URLs for testing
const indexerURL = "https://dev-indexer.api.gnoswap.io/graphql/query"

var (
	ErrFailedAllEndpoints = errors.New("failed to query from all endpoints")
)
