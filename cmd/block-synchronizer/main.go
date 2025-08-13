package main

import (
	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/graphql"
)

func main() {
	cfg := config.Load()
	_ = database.Connect(cfg.DB)
	_ = graphql.NewClient(cfg.GraphQL)
}
