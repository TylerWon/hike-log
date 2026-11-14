package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/router"
	"github.com/TylerWon/hike-log/backend/s3client"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/gin-gonic/gin"
)

const (
	SERVER_PORT           = ":8080"
	SHUTDOWN_TIMEOUT_SECS = 15
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.New()
	if err != nil {
		fatal("Failed to load config", err)
	}

	if cfg.Env != "local" {
		gin.SetMode(gin.ReleaseMode)
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

	handler, err := handler.NewHandler(db, store, photoStorage, auth, validatorManager)
	if err != nil {
		fatal("Failed to setup handler", err)
	}

	middleware := middleware.New(auth)

	router, err := router.New(handler, middleware, cfg.FrontendURL)
	if err != nil {
		fatal("Failed to setup router", err)
	}

	server := &http.Server{
		Addr:    SERVER_PORT,
		Handler: router,
	}
	go runServer(server)

	gracefulShutdown(server, db)
}

// Starts the HTTP server.
func runServer(server *http.Server) {
	// Blocks until error or shutdown
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fatal("Server error", err)
	}
}

// Handles shutdown of the HTTP server.
func gracefulShutdown(server *http.Server, db database.Database) {
	// Wait for termination signals (SIGTERM, SIGINT)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), SHUTDOWN_TIMEOUT_SECS*time.Second)
	defer cancel()

	// Shutdown server with timeout
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Shutdown error", "err", err)
	}

	// Close connection to database
	err = db.Close()
	if err != nil {
		slog.Error("Failed to close database connection", "err", err)
	}
}

// Logs a message then terminates the application.
func fatal(msg string, err error) {
	slog.Error(msg, "err", err)
	os.Exit(1)
}
