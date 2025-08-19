package main

import (
	"log/slog"

	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/logging"
	api "github.com/yanmoyy/tbi/internal/service/balance-api"
)

func main() {
	log := logging.NewLogger("balance-api")

	cfg := config.Load()
	db := database.NewClient(cfg.DB)
	service := api.New(cfg.API, db)
	log.Start()
	err := service.Run()
	if err != nil {
		slog.Error("Run", "err", err)
	}
	log.Finish()
}
