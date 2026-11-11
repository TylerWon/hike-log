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

	authHandler, err := handler.NewAuthHandler(auth)
	if err != nil {
		t.Fatal("Failed to setup auth handler: ", err)
	}

	return authHandler
}

// Creates a HikesHandler to use during testing.
func NewHikesHandler(t *testing.T, store store.Store, photoStorage photostorage.PhotoStorage) *handler.HikesHandler {
	t.Helper()

	hikesHandler, err := handler.NewHikesHandler(store, photoStorage)
	if err != nil {
		t.Fatal("Failed to setup hikes handler: ", err)
	}

	return hikesHandler
}

// Creates a PhotosHandler to use during testing.
func NewPhotosHandler(t *testing.T, store store.Store, photoStorage photostorage.PhotoStorage) *handler.PhotosHandler {
	t.Helper()

	photosHandler, err := handler.NewPhotosHandler(store, photoStorage)
	if err != nil {
		t.Fatal("Failed to setup photos handler: ", err)
	}

	return photosHandler
}
