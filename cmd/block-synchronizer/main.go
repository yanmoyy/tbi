package main

import (
	"log/slog"

	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/db"
)

func main() {
	cfg := config.Load()
	dbConn := db.Connect(cfg.DB)
	slog.Info("connected to db", "db", cfg.DB, "conn", dbConn)
}
