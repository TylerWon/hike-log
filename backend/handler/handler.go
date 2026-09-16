package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/TylerWon/hike-log/backend/aws"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// A Handler handles the request-response lifecycle for API routes.
type Handler struct {
	store    store.Store
	s3Client aws.S3Client
}

// Creates a new Handler
func New(store store.Store, s3Client aws.S3Client) *Handler {
	_, err := registerCustomValidators()
	if err != nil {
		log.Fatal("Error while registering custom validators: ", err)
	}

	return &Handler{store, s3Client}
}

// Returns a 200 OK. Used for application health checks.
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "ok"})
}

// Handles error that occurs when a record unexpectedly does not exist.
func handleRecordDoesNotExistError(c *gin.Context, err error, id uint) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": http.StatusText(http.StatusNotFound)})
		return
	}
	log.Printf("Failed to get record (id=%d): %v", id, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
}
