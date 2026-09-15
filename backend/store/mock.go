package store

import "github.com/TylerWon/hike-log/backend/models"

type MockStore struct {
	CreateRecordError   error
	DeleteRecordError   error
	GetHikeByIDResult  *models.Hike
	GetHikeByIDError   error
	GetPhotoByIDResult *models.Photo
	GetPhotoByIDError  error
	ListHikesResult    []models.Hike
	ListHikesError     error
	UpdateRecordError   error
}

func (m *MockStore) CloseConnection() error {
	return nil
}

func (m *MockStore) CreateRecord(mode interface{}) error {
	return m.CreateRecordError
}

func (m *MockStore) DeleteRecord(model interface{}) error {
	return m.DeleteRecordError
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

func (m *MockStore) UpdateRecord(model interface{}) error {
	return m.UpdateRecordError
}
