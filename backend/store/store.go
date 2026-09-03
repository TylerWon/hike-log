package store

import (
	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/models"
	"gorm.io/gorm"
)

// Store handles all interactions with the database for the app.
type Store interface {
	CloseConnection() error
	CreateHike(hike *models.Hike) error
	CreateHikes(hike []models.Hike) error
	GetHikeByID(id uint) (*models.Hike, error)
	ListHikes() ([]models.Hike, error)
}

// StoreImpl is an implementation of the Store interface.
type StoreImpl struct {
	db *gorm.DB
}

// Creates a Store connected to the database specified in the given dbConfig.
func New(dbConfig database.DbConfig) (Store, error) {
	db, err := database.Connect(dbConfig)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Hike{}, &models.Photo{})
	if err != nil {
		return nil, err
	}

	return &StoreImpl{db}, nil
}

// Ends the connection to the database.
func (store *StoreImpl) CloseConnection() error {
	sqlDb, err := store.db.DB()
	if err != nil {
		return err
	}

	err = sqlDb.Close()
	if err != nil {
		return err
	}

	return nil
}

// Creates a Hike.
func (store *StoreImpl) CreateHike(hike *models.Hike) error {
	result := store.db.Create(hike)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Creates Hikes.
func (store *StoreImpl) CreateHikes(hike []models.Hike) error {
	result := store.db.Create(hike)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Returns the Hike with the given ID.
func (store *StoreImpl) GetHikeByID(id uint) (*models.Hike, error) {
	var hike models.Hike

	result := store.db.First(&hike, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &hike, nil
}

// Returns all Hikes in reverse chronological order by Date.
func (store *StoreImpl) ListHikes() ([]models.Hike, error) {
	var hikes []models.Hike

	result := store.db.Preload("Photos").Order("date desc, trail_name").Find(&hikes)
	if result.Error != nil {
		return nil, result.Error
	}

	return hikes, nil
}
