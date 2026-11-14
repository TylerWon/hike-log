package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler handles the request-response lifecycle for API routes.
type Handler struct {
	Auth   *AuthHandler
	Health *HealthHandler
	Hikes  *HikesHandler
	Photos *PhotosHandler
}

// Creates a new Handler.
func NewHandler(
	db database.Database,
	store store.Store,
	photoStorage photostorage.PhotoStorage,
	auth auth.Auth,
	validatorManager ValidatorManager,
) (*Handler, error) {
	err := validatorManager.RegisterValidators()
	if err != nil {
		return nil, fmt.Errorf("Failed to register custom validators: %w", err)
	}

	return &Handler{
		Auth:   NewAuthHandler(auth),
		Health: NewHealthHandler(db),
		Hikes:  NewHikesHandler(store, photoStorage),
		Photos: NewPhotosHandler(store, photoStorage),
	}, nil
}

// Handles error that occurs when a record unexpectedly does not exist.
func handleRecordDoesNotExistError(c *gin.Context, err error, id uint) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": http.StatusText(http.StatusNotFound)})
		return
	}
	slog.Error("Failed to get record", "id", id, "err", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
}
