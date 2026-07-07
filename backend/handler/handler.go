package handler

import (
	"errors"
	"fmt"
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

// Returns all Hikes in reverse chronological order by Date.
// For each Hike, its CoverPhoto is included in the response, but not its Photos.
func (h *Handler) ListHike(c *gin.Context) {
	var hikes []models.Hike

	result := h.db.Preload("CoverPhoto").Order("date desc, trail_name").Find(&hikes)
	if result.Error != nil {
		log.Println("Failed to list tasks: ", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, hikes)
}

// Returns the Hike with an ID matching the id path parameter. If no such Hike exists, a 404 error is returned instead.
// For each Hike, its Photos are included in the response, but not its CoverPhoto.
func (h *Handler) RetrieveHike(c *gin.Context) {
	hikeId := c.Param("id")
	var hike models.Hike

	result := h.db.Preload("Photos").First(&hike, hikeId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Hike %v not found", hikeId)})
			return
		}
		log.Printf("Failed to retrieve hike %v: %v", hikeId, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, hike)
}
