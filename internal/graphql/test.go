package graphql

import (
	"errors"
	"flag"
)

// Endpoint URLs for testing
const indexerURL = "https://dev-indexer.api.gnoswap.io/graphql/query"

var (
	ErrFailedAllEndpoints = errors.New("failed to query from all endpoints")
)

// Flags for testing
var offline = flag.Bool("offline", false, "Run offline")
var minBlocks = flag.Int("min-blocks", 0, "Minimum number of blocks to receive before exiting")
