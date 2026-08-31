package handler

type createPhotoUploadURLPathParams struct {
	// The ID of the Hike the photo belongs to.
	HikeID uint `uri:"hikeId" binding:"required"`
}
