package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
