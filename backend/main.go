package main

import (
	"log"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/router"
	"github.com/TylerWon/hike-log/backend/s3"
	"github.com/TylerWon/hike-log/backend/store"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	db, err := database.Connect(config.Database)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	store, err := store.New(db)
	if err != nil {
		log.Fatal("Failed to create store: ", err)
	}

	s3Client, err := s3.NewS3Client(config.Env, config.AWS)
	if err != nil {
		log.Fatal("Failed to setup S3 client: ", err)
	}

	photoStorage, err := photostorage.New(s3Client, config.Env, config.AWS)
	if err != nil {
		log.Fatal("Failed to setup photo storage: ", err)
	}

	auth, err := auth.New(config.Auth)
	if err != nil {
		log.Fatal("Failed to setup auth: ", err)
	}

	validatorManager, err := handler.NewValidatorManager()
	if err != nil {
		log.Fatal("Failed to setup validator manager: ", err)
	}

	handler, err := handler.NewHandler(store, photoStorage, auth, validatorManager)
	if err != nil {
		log.Fatal("Failed to setup handler: ", err)
	}

	middleware := middleware.New(auth)

	router, err := router.New(handler, middleware, config.FrontendURL)
	if err != nil {
		log.Fatal("Failed to setup router: ", err)
	}

	router.Run()
}
