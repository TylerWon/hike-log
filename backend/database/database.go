package database

import (
	"fmt"
	"os"

	"github.com/TylerWon/todo-app/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect to the database and migrate the models
func Connect() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", dbHost, dbPort, dbName, dbUser, dbPassword)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Hike{})
	db.AutoMigrate(&models.Photo{})

	return db, nil
}
