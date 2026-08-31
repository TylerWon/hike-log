package handler

type createHikeRequest struct {
	// The name of the trail that was hiked.
	TrailName string `json:"trailName" binding:"required"`

	// The date of the hike in YYYY-MM-DD format.
	Date string `json:"date" binding:"required,datetime=2006-01-02"`

	// Notes on the hike.
	Notes string `json:"notes" binding:"required"`

	// How enjoyable the hike was out of 5. Only whole and half values are allowed.
	Rating float32 `json:"rating" binding:"required,gte=0,lte=5,divisibleByHalf"`

	// How difficult the hike was out of 10. Only whole and half values are allowed.
	Difficulty float32 `json:"difficulty" binding:"required,gte=0,lte=10,divisibleByHalf"`

	// The distance hiked (km).
	Distance float32 `json:"distance" binding:"required"`

	// The elevation gained on the hike (m).
	ElevationGain uint `json:"elevationGain" binding:"required"`

	// Time it took to complete the hike (mins).
	Duration uint `json:"duration" binding:"required"`

	// Link to the AllTrails page for the trail.
	AllTrailsUrl string `json:"allTrailsUrl" binding:"required,url"`
}

type createPhotoUploadURLRequest struct {
	// MIME type of the image. Must be a valid image type.
	ContentType string `json:"contentType" binding:"required,validImageType"`

	// Size of the image in bytes. Must be between 0 to 10 MB.
	ContentLength uint `json:"contentLength" binding:"required,gte=0,lte=10485760"`
}
