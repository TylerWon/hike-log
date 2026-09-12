package store

import "github.com/TylerWon/hike-log/backend/models"

type MockStore struct {
	CreateModelError  error
	GetHikeByIDResult *models.Hike
	GetHikeByIDError  error
	ListHikesResult   []models.Hike
	ListHikesError    error
}

func (m *MockStore) CloseConnection() error {
	return nil
}

func (m *MockStore) CreateModel(mode interface{}) error {
	return m.CreateModelError
}

func (m *MockStore) GetHikeByID(id uint) (*models.Hike, error) {
	return m.GetHikeByIDResult, m.GetHikeByIDError
}

func (m *MockStore) ListHikes() ([]models.Hike, error) {
	return m.ListHikesResult, m.ListHikesError
}
