package test

import (
	"flag"
	"testing"
)

var testDB = flag.Bool("test.db", false, "run database tests")
var testIndexer = flag.Bool("test.indexer", false, "run indexer tests")

// WARNING: DB tests should be run with testing environment
func CheckDBFlag(t *testing.T) {
	flag.Parse()
	t.Helper()
	if !*testDB {
		t.Skip("Skipping test in non-dbtest mode")
	}
}

func CheckIndexerFlag(t *testing.T) {
	flag.Parse()
	t.Helper()
	if !*testIndexer {
		t.Skip("Skipping test in non-indexer mode")
	}
}
