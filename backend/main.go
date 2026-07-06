package main

import (
	"log"
	"os"

	"github.com/TylerWon/todo-app/backend/database"
	"github.com/TylerWon/todo-app/backend/handler"
	"github.com/TylerWon/todo-app/backend/router"
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

	handler := handler.New(db)
	router := router.New(handler)

	router.Run()
}
