package testutils

import (
	"fmt"
	"testing"
	"time"

	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/store"
	"gorm.io/datatypes"
)

// Creates n Hikes and returns them. Saves them to the database if specified.
func CreateHikes(t *testing.T, store store.Store, n int, save bool) []models.Hike {
	var hikes []models.Hike
	for i := range n {
		photos := []models.Photo{
			{SrcUrl: "https://example.com/photo-1.jpg"},
		}
		hike := models.Hike{
			TrailName:     fmt.Sprintf("Trail %d", i),
			Date:          datatypes.Date(time.Date(2026, 1, i, 0, 0, 0, 0, time.UTC)),
			Notes:         "Hike notes",
			Rating:        3,
			Difficulty:    9,
			Distance:      10,
			ElevationGain: 1000,
			Duration:      120,
			AllTrailsUrl:  "https://www.alltrails.com/",
			Photos:        photos,
		}
		hikes = append(hikes, hike)
	}

	if save {
		err := store.CreateHikes(hikes)
		if err != nil {
			t.Fatal("Failed to save hikes: ", err)
		}
	}

	return hikes
}
