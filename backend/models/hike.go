package models

import (
	"gorm.io/datatypes"
)

type Hike struct {
	// Primary key
	ID uint `json:"id"`

	// The trail that was hiked.
	Trail string `json:"trail"`

	// The date of the hike.
	Date datatypes.Date `json:"date"`

	// Notes on the hike.
	Notes *string `json:"notes"`

	// How enjoyable the hike was out of 5. Half values allowed.
	Rating float32 `json:"rating"`

	// How difficult the hike was out of 10. Half values allowed.
	Difficulty float32 `json:"difficulty"`

	// The distance (km) hiked.
	Distance float32 `json:"distance"`

	// The elevation (m) gained on the hike.
	ElevationGain uint `json:"elevationGain"`

	// Link to the AllTrails page for the trail.
	AllTrailsUrl *string `json:"allTrailsUrl"`

	// Photos taken on the hike.
	Photos []Photo `json:"photos"`

	// Photo to use as the thumbnail for the hike.
	CoverPhotoID *uint  `json:"coverPhotoId"`
	CoverPhoto   *Photo `json:"coverPhoto"`
}
