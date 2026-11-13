package main

import (
	"log/slog"
	"os"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/router"
	"github.com/TylerWon/hike-log/backend/s3client"
	"github.com/TylerWon/hike-log/backend/store"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.New()
	if err != nil {
		fatal("Failed to load config", err)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		fatal("Failed to connect to database", err)
	}

	store, err := store.New(db)
	if err != nil {
		fatal("Failed to create store", err)
	}

	s3Client, err := s3client.New(cfg.Env, cfg.AWS)
	if err != nil {
		fatal("Failed to setup S3 client", err)
	}

	photoStorage, err := photostorage.New(s3Client, cfg.Env, cfg.AWS)
	if err != nil {
		fatal("Failed to setup photo storage", err)
	}

	auth, err := auth.New(cfg.Auth)
	if err != nil {
		fatal("Failed to setup auth", err)
	}

	validatorManager, err := handler.NewValidatorManager()
	if err != nil {
		fatal("Failed to setup validator manager", err)
	}

	handler, err := handler.NewHandler(store, photoStorage, auth, validatorManager)
	if err != nil {
		fatal("Failed to setup handler", err)
	}

	middleware := middleware.New(auth)

	router, err := router.New(handler, middleware, cfg.FrontendURL)
	if err != nil {
		fatal("Failed to setup router", err)
	}

	if err := router.Run(); err != nil {
		fatal("Failed to start server", err)
	}
}

func fatal(msg string, err error) {
	slog.Error(msg, "err", err)
	os.Exit(1)
}
