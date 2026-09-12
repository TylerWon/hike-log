package handler

type CreatePhotoUploadURLResponse struct {
	// The presigned URL that can be used to upload a photo to the S3 bucket.
	UploadURL string `json:"uploadUrl"`

	// The key that is assigned to the photo when it is uploaded to the S3 bucket. Format "hikes/<hike_id>/photos/<uuid>".
	ObjectKey string `json:"objectKey"`
}
