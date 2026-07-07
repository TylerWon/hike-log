package models

import (
	"gorm.io/datatypes"
)

// A Hike represents a hike that was completed.
type Hike struct {
	// Primary key
	ID uint `json:"id"`

	// The name of the trail that was hiked.
	TrailName string `json:"trail"`

	// The date of the hike.
	Date datatypes.Date `json:"date"`

	// Notes on the hike.
	Notes string `json:"notes"`

	// How enjoyable the hike was out of 5. Half values allowed.
	Rating float32 `json:"rating"`

	// How difficult the hike was out of 10. Half values allowed.
	Difficulty float32 `json:"difficulty"`

	// The distance hiked (km).
	Distance float32 `json:"distance"`

	// The elevation gained on the hike (m).
	ElevationGain uint `json:"elevationGain"`

	// Time it took to complete the hike (mins).
	TotalTime uint `json:"totalTime"`

	// Link to the AllTrails page for the trail.
	AllTrailsUrl string `json:"allTrailsUrl"`

	// Photos taken on the hike.
	Photos []Photo `json:"photos"`

	// Photo to use as the thumbnail for the hike.
	CoverPhotoID *uint  `json:"-"` // Omit when struct is serialized into JSON
	CoverPhoto   *Photo `json:"coverPhoto"`
}
