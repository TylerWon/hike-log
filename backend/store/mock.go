package store

import "github.com/TylerWon/hike-log/backend/models"

type MockStore struct {
	CreateModelError   error
	DeleteModelError   error
	GetHikeByIDResult  *models.Hike
	GetHikeByIDError   error
	GetPhotoByIDResult *models.Photo
	GetPhotoByIDError  error
	ListHikesResult    []models.Hike
	ListHikesError     error
	UpdateModelError   error
}

func (m *MockStore) CloseConnection() error {
	return nil
}

func (m *MockStore) CreateModel(mode interface{}) error {
	return m.CreateModelError
}

func (m *MockStore) DeleteModel(model interface{}) error {
	return m.DeleteModelError
}

func (m *MockStore) GetHikeByID(id uint) (*models.Hike, error) {
	return m.GetHikeByIDResult, m.GetHikeByIDError
}

func (m *MockStore) GetPhotoByID(id uint) (*models.Photo, error) {
	return m.GetPhotoByIDResult, m.GetPhotoByIDError
}

func (m *MockStore) ListHikes() ([]models.Hike, error) {
	return m.ListHikesResult, m.ListHikesError
}

func (m *MockStore) UpdateModel(model interface{}) error {
	return m.UpdateModelError
}
