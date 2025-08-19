package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondServerError(c *gin.Context, respondMsg string, err error) {
	if err != nil {
		slog.Error("respondServerError", "error", err, "respondMsg", respondMsg)
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": respondMsg})
}

func respondNotFound(c *gin.Context, respondMsg string) {
	c.JSON(http.StatusNotFound, gin.H{"error": respondMsg})
}
