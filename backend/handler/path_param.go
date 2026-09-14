package handler

type hikeIDPathParam struct {
	// The ID of the Hike.
	HikeID uint `uri:"hikeId" binding:"required"`
}
