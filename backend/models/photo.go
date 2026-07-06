package models

// A Photo represents a photo taken on a Hike.
type Photo struct {
	// Primary key
	ID uint `json:"id"`

	// Link to where the photo can be accessed.
	SrcUrl string `json:"srcUrl"`

	// Caption for the photo.
	Caption *string `json:"caption"`

	// The Hike the photo was taken on.
	HikeID uint `json:"hikeId"`
}
