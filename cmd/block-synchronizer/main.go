package main

import (
	"fmt"
	"log/slog"

	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/indexer"
	synchronizer "github.com/yanmoyy/tbi/internal/service/block-synchronizer"
)

func main() {
	fmt.Println()
	slog.Info("##### block-synchronizer #####")
	fmt.Println()
	cfg := config.Load()
	db := database.NewClient(cfg.DB)
	graphql := indexer.NewClient(cfg.GraphQL)
	s := synchronizer.New(graphql, db, nil)
	slog.Info("block-synchronizer Initialized!")
	slog.Info("Starting backfill...")
	s.RunBackfill()
	slog.Info("Backfill finished!")
}
