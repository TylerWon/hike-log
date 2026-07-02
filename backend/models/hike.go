package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Hike struct {
	gorm.Model
	Trail        string
	Date         datatypes.Date
	Notes        *string
	Rating       float32
	Difficulty   float32
	Distance     float32
	Elevation    uint
	AllTrailsUrl *string
	Photos       []Photo
}
