package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/database"
)

type Service struct {
	db     *database.Client
	router *gin.Engine
	port   string
}

func New(cfg config.API, db *database.Client) *Service {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	s := &Service{
		db:     db,
		router: router,
		port:   cfg.Port,
	}
	s.setupRoutes()
	return s
}

func (s *Service) Run() error {
	return s.router.Run(":" + s.port)
}

func (s *Service) setupRoutes() {
	s.router.GET("/tokens/balances", s.handleGetTokenBalances)
	s.router.GET("/tokens/:tokenPath/balances", s.handleGetTokenAccountBalances)
	s.router.GET("/tokens/transfer-history", s.handleGetTransferHistory)
	// check health
	s.router.GET("/health", s.handleHealthCheck)
}
