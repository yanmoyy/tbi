package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
	ctx := context.Background()

	slog.Info("block-synchronizer Started...")
	defer slog.Info("block-synchronizer Finished!")
	err := s.Run(ctx)
	if err != nil {
		slog.Error("Run", "err", err)
		return
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
