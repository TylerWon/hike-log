package main

import (
	"log"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/config"
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

	store, err := store.New(config.Database)
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

	middleware := middleware.New(auth)
	handler := handler.NewHandler(store, photoStorage, auth)

	router, err := router.New(handler, middleware, config.FrontendURL)
	if err != nil {
		log.Fatal("Failed to setup router: ", err)
	}

	router.Run()
}
