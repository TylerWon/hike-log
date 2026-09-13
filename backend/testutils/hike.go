package testutils

import (
	"fmt"
	"testing"
	"time"

	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/store"
	"gorm.io/datatypes"
)

// Constructs n Hikes and returns them. Optionally can add Photos to the Hikes and save the Hikes to the database.
func ConstructHikes(t *testing.T, n int, store store.Store, photos bool, save bool) []models.Hike {
	t.Helper()

	var hikes []models.Hike
	for i := range n {
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
		}

		if photos {
			hike.Photos = []models.Photo{
				{
					SrcUrl:       fmt.Sprintf("http://localstack:4566/hike-log/hikes/%d/photos/acde070d-8c4c-4f0d-9d8a-162843c10333", i),
					Caption:      "Caption",
					DisplayOrder: 1,
				},
			}
		}

		hikes = append(hikes, hike)
	}

	if save {
		err := store.CreateModel(hikes)
		if err != nil {
			t.Fatal("Failed to save hikes: ", err)
		}
	}

	return hikes
}
