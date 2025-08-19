package test

import (
	"flag"
	"testing"
)

var testDB = flag.Bool("test.db", false, "run database tests")
var testSQS = flag.Bool("test.sqs", false, "run SQS tests")
var testAPI = flag.Bool("test.api", false, "run API tests")
var testIndexer = flag.Bool("test.indexer", false, "run indexer tests")

// Just for importing test package to avoid flag.Parse() in every test
func NoFlag(t *testing.T) {
	t.Helper()
	flag.Parse()
}

func CheckAPIFlag(t *testing.T) {
	t.Helper()
	flag.Parse()
	if !*testAPI {
		t.Skip("Skipping test in non-api mode")
	}
}

// WARNING: DB tests should be run with testing environment
func CheckDBFlag(t *testing.T) {
	t.Helper()
	flag.Parse()
	if !*testDB {
		t.Skip("Skipping test in non-db mode")
	}
}

// WARNING: SQS tests should be run with testing environment
func CheckSQSFlag(t *testing.T) {
	flag.Parse()
	t.Helper()
	if !*testSQS {
		t.Skip("Skipping test in non-sqs mode")
	}
}

func CheckIndexerFlag(t *testing.T) {
	t.Helper()
	flag.Parse()
	if !*testIndexer {
		t.Skip("Skipping test in non-indexer mode")
	}
}
