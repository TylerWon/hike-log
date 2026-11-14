package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/TylerWon/hike-log/backend/database"
	"github.com/gin-gonic/gin"
)

const readyCheckTimeout = 2 * time.Second

// HealthHandler handles health and readiness probe routes.
type HealthHandler struct {
	db database.Database
}

// Creates a new HealthHandler.
func NewHealthHandler(db database.Database) *HealthHandler {
	return &HealthHandler{db: db}
}

// Returns a 200 OK for application liveness checks.
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "ok"})
}

// Returns 200 when the database is reachable, otherwise 503.
func (h *HealthHandler) ReadyCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), readyCheckTimeout)
	defer cancel()

	err := h.db.Ping(ctx)
	if err != nil {
		slog.Error("Readiness check failed", "err", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"ready": "unavailable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ready": "ok"})
}
