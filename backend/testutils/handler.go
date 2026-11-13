package testutils

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/store"
)

// Creates an AuthHandler to use during testing.
func NewAuthHandler(t *testing.T, auth auth.Auth) *handler.AuthHandler {
	t.Helper()

	_, err := handler.RegisterCustomValidators()
	if err != nil {
		t.Fatal("Failed to register custom validators: ", err)
	}

	return handler.NewAuthHandler(auth)
}

// Creates a HikesHandler to use during testing.
func NewHikesHandler(t *testing.T, store store.Store, photoStorage photostorage.PhotoStorage) *handler.HikesHandler {
	t.Helper()

	_, err := handler.RegisterCustomValidators()
	if err != nil {
		t.Fatal("Failed to register custom validators: ", err)
	}

	return handler.NewHikesHandler(store, photoStorage)
}

// Creates a PhotosHandler to use during testing.
func NewPhotosHandler(t *testing.T, store store.Store, photoStorage photostorage.PhotoStorage) *handler.PhotosHandler {
	t.Helper()

	_, err := handler.RegisterCustomValidators()
	if err != nil {
		t.Fatal("Failed to register custom validators: ", err)
	}

	return handler.NewPhotosHandler(store, photoStorage)
}
