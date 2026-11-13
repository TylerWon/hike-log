package testutils

import (
	"fmt"
	"testing"
	"time"

	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/models/types"
	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/store"
)

// Constructs n Hikes and returns them. Optionally allows a Photo to be created for each Hike.
func ConstructHikes(t *testing.T, n int, photoStorage photostorage.PhotoStorage, photo bool) []models.Hike {
	t.Helper()

	var hikes []models.Hike
	for i := range n {
		hike := models.Hike{
			TrailName:     fmt.Sprintf("Trail %d", i),
			Date:          types.Date(time.Date(2026, 1, i+1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")),
			Notes:         "Hike notes",
			Rating:        3,
			Difficulty:    9,
			Distance:      10,
			ElevationGain: 1000,
			Duration:      120,
			AllTrailsUrl:  "https://www.alltrails.com/",
		}

		if photo {
			objectKey := photoStorage.CreatePhotoObjectKey(uint(i))
			hike.Photos = []models.Photo{
				{
					SrcUrl:       photoStorage.GetPhotoObjectURL(objectKey),
					Caption:      "Caption",
					DisplayOrder: 1,
				},
			}
		}

		hikes = append(hikes, hike)
	}

	return hikes
}

// Constructs n Hikes and saves them to the database. Optionally allows a Photo to be created for each Hike.
func ConstructAndSaveHikes(
	t *testing.T,
	n int,
	store store.Store,
	photoStorage photostorage.PhotoStorage,
	photo bool,
) []models.Hike {
	t.Helper()

	hikes := ConstructHikes(t, n, photoStorage, photo)

	err := store.CreateRecord(hikes)
	if err != nil {
		t.Fatal("Failed to save hikes: ", err)
	}

	return hikes
}
