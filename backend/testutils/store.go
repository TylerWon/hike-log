package testutils

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/store"
)

// Creates a Store backed by the given database.
func NewStore(t *testing.T, db database.Database) store.Store {
	t.Helper()

	s, err := store.New(db)
	if err != nil {
		t.Fatal("Failed to create store: ", err)
	}

	return s
}
