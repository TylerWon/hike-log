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
	authHandler, err := NewAuthHandler(auth)
	if err != nil {
		return nil, fmt.Errorf("Failed to create auth handler: %w", err)
	}

	hikesHandler, err := NewHikesHandler(store, photoStorage)
	if err != nil {
		return nil, fmt.Errorf("Failed to create hikes handler: %w", err)
	}

	photosHandler, err := NewPhotosHandler(store, photoStorage)
	if err != nil {
		return nil, fmt.Errorf("Failed to create photos handler: %w", err)
	}

	return &Handler{
		Auth:   authHandler,
		Hikes:  hikesHandler,
		Photos: photosHandler,
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
