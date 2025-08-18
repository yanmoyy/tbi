package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/indexer"
	"github.com/yanmoyy/tbi/internal/logging"
	synchronizer "github.com/yanmoyy/tbi/internal/service/block-synchronizer"
	"github.com/yanmoyy/tbi/internal/sqs"
)

func main() {
	log := logging.NewLogger("block-synchronizer")
	cfg := config.Load()

	db := database.NewClient(cfg.DB)
	graphql := indexer.NewClient(cfg.GraphQL)

	ctx := context.Background()
	sqs := sqs.NewClient(cfg.SQS)

	service := synchronizer.New(graphql, db, sqs)

	log.Start()

	err := service.RunBackFill(ctx)
	if err != nil {
		slog.Error("RunBackFill", "err", err)
		return
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Finish()
}
