package handler

type createPhotoUploadURLPathParams struct {
	// The ID of the Hike the photo belongs to.
	HikeID uint `uri:"hikeId" binding:"required"`
}

type createPhotoPathParams struct {
	// The ID of the Hike the photo belongs to.
	HikeID uint `uri:"hikeId" binding:"required"`
}

type deleteHikePathParams struct {
	// The ID of the Hike to delete.
	HikeID uint `uri:"hikeId" binding:"required"`
}
