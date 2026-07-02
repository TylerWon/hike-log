package main

import (
	"log"
	"net/http"

	"github.com/TylerWon/todo-app/backend/database"
	"github.com/gin-gonic/gin"
)

func main() {
	_, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.Run()
}
