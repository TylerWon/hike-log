package handler

import "github.com/TylerWon/hike-log/backend/models"

type createPhotosResponseItem struct {
	// Outcome for this item.
	Success bool `json:"success"`

	// The newly created Photo when "Success" is true.
	Result *models.Photo `json:"result,omitempty"`

	// Error message when "Success" is false.
	Error string `json:"error,omitempty"`
}

type CreatePhotosResponse []createPhotosResponseItem

type createPresignedURLsResponseItemResult struct {
	// The presigned URL that can be used to upload a photo to the S3 bucket.
	PresignedURL string `json:"presignedUrl"`

	// The key that is assigned to the photo when it is uploaded to the S3 bucket. Format "hikes/<hike_id>/photos/<uuid>".
	ObjectKey string `json:"objectKey"`
}

type createPresignedURLsResponseItem struct {
	// Outcome for this item.
	Success bool `json:"success"`

	// The presigned URL data when "Success" is true.
	Result *createPresignedURLsResponseItemResult `json:"result,omitempty"`

	// Error message when "Success" is false.
	Error string `json:"error,omitempty"`
}

type CreatePresignedURLsResponse []createPresignedURLsResponseItem
