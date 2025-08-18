package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/logging"
	processor "github.com/yanmoyy/tbi/internal/service/event_processor"
	"github.com/yanmoyy/tbi/internal/sqs"
)

func main() {
	log := logging.NewLogger("event-processor")

	cfg := config.Load()
	db := database.NewClient(cfg.DB)

	sqs := sqs.NewClient(cfg.SQS)

	service := processor.New(db, sqs)

	log.Start()
	ctx := context.Background()

	err := service.Run(ctx)
	if err != nil {
		slog.Error("Run", "err", err)
		return
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Finish()
}
