package main

import (
	"log/slog"

	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/indexer"
)

func main() {
	cfg := config.Load()
	_ = database.NewClient(cfg.DB)
	_ = indexer.NewClient(cfg.GraphQL)
	slog.Info("block-synchronizer started")
}
