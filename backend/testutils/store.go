package testutils

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/store"
)

// Creates a Store connected to the test suite database.
func NewStore(t *testing.T, testSuiteDB *TestSuiteDB) store.Store {
	t.Helper()

	store, err := store.New(testSuiteDB.testDBConfig)
	if err != nil {
		t.Fatal("Failed to create store: ", err)
	}

	return store
}

// Closes the database connection held by the Store.
func TeardownStore(t *testing.T, store store.Store) {
	t.Helper()

	err := store.CloseConnection()
	if err != nil {
		t.Fatal("Failed to close store database connection: ", err)
	}
}
