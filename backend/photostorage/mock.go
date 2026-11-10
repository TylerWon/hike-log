package photostorage

import "context"

type MockPhotoStorage struct {
	CreatePhotoObjectKeyResult     string
	CreatePresignedUploadURLResult string
	CreatePresignedUploadURLKey    string
	CreatePresignedUploadURLError  error
	DeletePhotoError               error
	DeletePhotosForHikeError       error
	DoesPhotoExistResult           bool
	DoesPhotoExistError            error
	GetPhotoObjectKeyResult        string
	GetPhotoObjectURLResult        string
}

func (m *MockPhotoStorage) CreatePhotoObjectKey(hikeID uint) string {
	return m.CreatePhotoObjectKeyResult
}

func (m *MockPhotoStorage) CreatePresignedUploadURL(
	ctx context.Context,
	hikeID uint,
	contentType string,
	contentLength int64,
) (string, string, error) {
	return m.CreatePresignedUploadURLResult, m.CreatePresignedUploadURLKey, m.CreatePresignedUploadURLError
}

func (m *MockPhotoStorage) DeletePhoto(ctx context.Context, objectURL string) error {
	return m.DeletePhotoError
}

func (m *MockPhotoStorage) DeletePhotosForHike(ctx context.Context, hikeID uint) error {
	return m.DeletePhotosForHikeError
}

func (m *MockPhotoStorage) DoesPhotoExist(ctx context.Context, objectKey string) (bool, error) {
	return m.DoesPhotoExistResult, m.DoesPhotoExistError
}

func (m *MockPhotoStorage) GetPhotoObjectKey(objectURL string) string {
	return m.GetPhotoObjectKeyResult
}

func (m *MockPhotoStorage) GetPhotoObjectURL(objectKey string) string {
	return m.GetPhotoObjectURLResult
}
