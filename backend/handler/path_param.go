package handler

type hikeIDPathParam struct {
	// The ID of the Hike.
	HikeID uint `uri:"hikeId" binding:"required"`
}

type photoIDPathParam struct {
	// The ID of the Photo.
	PhotoID uint `uri:"photoId" binding:"required"`
}
