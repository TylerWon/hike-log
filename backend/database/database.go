package database

import (
	"fmt"

	"github.com/TylerWon/todo-app/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string
}

func Setup(dbConfig DbConfig) (*gorm.DB, error) {
	db, err := Connect(dbConfig)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Hike{}, &models.Photo{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Connect(dbConfig DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s",
		dbConfig.DbHost,
		dbConfig.DbPort,
		dbConfig.DbName,
		dbConfig.DbUser,
		dbConfig.DbPassword,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
