package main

import (
	"log"
	"os"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/router"
	"github.com/TylerWon/hike-log/backend/s3"
	"github.com/TylerWon/hike-log/backend/store"
)

func main() {
	dbConfig := database.DBConfig{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbName:     os.Getenv("DB_NAME"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
	}
	store, err := store.New(dbConfig)
	if err != nil {
		log.Fatal("Failed to create store: ", err)
	}

	s3Client, err := s3.NewS3Client()
	if err != nil {
		log.Fatal("Failed to setup S3 client: ", err)
	}

	auth, err := auth.New(os.Getenv("ADMIN_PASSWORD_HASH"), os.Getenv("SESSION_SECRET"))
	if err != nil {
		log.Fatal("Failed to setup auth: ", err)
	}

	middleware := middleware.New(auth)
	handler := handler.New(store, s3Client, auth)
	router := router.New(handler, middleware)

	router.Run()
}
