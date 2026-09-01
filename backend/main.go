package main

import (
	"log"
	"os"

	"github.com/TylerWon/hike-log/backend/aws"
	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/router"
)

func main() {
	dbConfig := database.DbConfig{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbName:     os.Getenv("DB_NAME"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
	}

	db, err := database.Setup(dbConfig)
	if err != nil {
		log.Fatal("Failed to setup database: ", err)
	}

	s3Client, err := aws.NewS3Client()
	if err != nil {
		log.Fatal("Failed to setup S3 client: ", err)
	}

	handler := handler.New(db, s3Client)
	router := router.New(handler)

	router.Run()
}
