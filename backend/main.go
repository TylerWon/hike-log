package main

import (
	"log"
	"os"

	"github.com/TylerWon/hike-log/backend/aws"
	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/router"
	"github.com/TylerWon/hike-log/backend/store"
)

func main() {
	dbConfig := database.DbConfig{
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

	s3Client, err := aws.NewS3Client()
	if err != nil {
		log.Fatal("Failed to setup S3 client: ", err)
	}

	handler := handler.New(store, s3Client)
	router := router.New(handler)

	router.Run()
}
