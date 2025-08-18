package database

import (
	"flag"
	"testing"

	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/test"
)

var dbHost = flag.String("db.host", "localhost", "Database host")

func getTestClient(t *testing.T) *Client {
	t.Helper()
	test.CheckDBFlag(t)
	cfg := config.LoadWithPath("../../.env")
	cfg.DB.Host = *dbHost
	c := NewClient(cfg.DB)
	return c
}
