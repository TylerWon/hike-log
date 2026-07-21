package handler

import (
	"log"
	"net/http"

	"github.com/TylerWon/hike-log/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// A Handler handles the request-response lifecycle for API routes.
type Handler struct {
	db *gorm.DB
}

// Creates a new Handler with access to the given database.
func New(db *gorm.DB) *Handler {
	return &Handler{db}
}

// Returns a 200 OK. Used for application health checks.
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "ok"})
}

// Returns all Hikes in reverse chronological order by Date.
func (h *Handler) ListHike(c *gin.Context) {
	var hikes []models.Hike

	result := h.db.Preload("Photos").Order("date desc, trail_name").Find(&hikes)
	if result.Error != nil {
		log.Println("Failed to list hikes: ", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, hikes)
}
