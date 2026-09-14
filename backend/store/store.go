package store

import (
	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/models"
	"gorm.io/gorm"
)

// Store handles all interactions with the database for the app.
type Store interface {
	CloseConnection() error
	CreateModel(model interface{}) error
	DeleteModel(model interface{}) error
	GetHikeByID(id uint) (*models.Hike, error)
	GetPhotoByID(id uint) (*models.Photo, error)
	ListHikes() ([]models.Hike, error)
	UpdateModel(model interface{}) error
}

// storeImpl is an implementation of the Store interface.
type storeImpl struct {
	db database.Database
}

// Creates a Store connected to the database specified in the given dbConfig.
func New(dbConfig database.DBConfig) (Store, error) {
	db, err := database.Connect(dbConfig)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Hike{}, &models.Photo{})
	if err != nil {
		return nil, err
	}

	return &storeImpl{db}, nil
}

// Creates a new Store. Allows injection of internal dependencies to allow for mocking.
func NewTestStore(db database.Database) Store {
	return &storeImpl{db}
}

// Ends the connection to the database.
func (store *storeImpl) CloseConnection() error {
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

// Creates the provided model/models.
func (store *storeImpl) CreateModel(model interface{}) error {
	result := store.db.Create(model)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Deletes the provided model/models. The model(s) should have their ID set. Deleting a non-existent model does not
// result in error.
func (store *storeImpl) DeleteModel(model interface{}) error {
	result := store.db.Delete(model)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Returns the Hike with the given ID.
func (store *storeImpl) GetHikeByID(id uint) (*models.Hike, error) {
	var hike models.Hike

	result := store.db.First(&hike, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &hike, nil
}

// Returns the Photo with the given ID.
func (store *storeImpl) GetPhotoByID(id uint) (*models.Photo, error) {
	var photo models.Photo

	result := store.db.First(&photo, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &photo, nil
}

// Returns all Hikes and their Photos in reverse chronological order by Date. The Photos are sorted by DisplayOrder.
func (store *storeImpl) ListHikes() ([]models.Hike, error) {
	var hikes []models.Hike

	result := store.db.
		Preload("Photos", func(db *gorm.DB) *gorm.DB { return db.Order("display_order") }).
		Order("date desc, trail_name").Find(&hikes)
	if result.Error != nil {
		return nil, result.Error
	}

	return hikes, nil
}

// Updates the provided model. Updating a non-existent model does not result in error.
func (store *storeImpl) UpdateModel(model interface{}) error {
	result := store.db.Update(model)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
