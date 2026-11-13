package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler handles the request-response lifecycle for API routes.
type Handler struct {
	Auth   *AuthHandler
	Hikes  *HikesHandler
	Photos *PhotosHandler
}

// Creates a new Handler.
func NewHandler(store store.Store, photoStorage photostorage.PhotoStorage, auth auth.Auth) (*Handler, error) {
	_, err := RegisterCustomValidators()
	if err != nil {
		return nil, fmt.Errorf("Failed to register custom validators: %w", err)
	}

	return &Handler{
		Auth:   NewAuthHandler(auth),
		Hikes:  NewHikesHandler(store, photoStorage),
		Photos: NewPhotosHandler(store, photoStorage),
	}, nil
}

// Returns a 200 OK for application health checks.
func HealthCheck(c *gin.Context) {
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
