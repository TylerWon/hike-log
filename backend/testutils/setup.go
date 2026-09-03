package testutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/router"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/gin-gonic/gin"
)

// Configures a fresh database to use during testing.
func SetupTestDB(t *testing.T) store.Store {
	// Connect to the default 'postgres' database
	dbConfig := database.DbConfig{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbName:     "postgres",
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
	}

	db, err := database.Connect(dbConfig)
	if err != nil {
		t.Fatal("Failed to connect to 'postgres' database: ", err)
	}

	// Create test database
	testDbName := "test_" + os.Getenv("DB_NAME")
	var exists int
	db.Raw("SELECT 1 FROM pg_database WHERE datname = ?", testDbName).Scan(&exists)

	if exists != 0 {
		result := db.Exec(fmt.Sprintf("DROP DATABASE %s", testDbName))
		if result.Error != nil {
			t.Fatal("Failed to delete old test database: ", result.Error)
		}
	}

	result := db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDbName))
	if result.Error != nil {
		t.Fatal("Failed to create test database: ", result.Error)
	}

	// Create store connected to the test database
	dbConfig.DbName = testDbName
	store, err := store.New(dbConfig)
	if err != nil {
		t.Fatal("Failed to setup test database: ", err)
	}

	return store
}

// Cleans up the test database.
func TeardownTestDB(t *testing.T, store store.Store) {
	err := store.CloseConnection()
	if err != nil {
		t.Fatal("Failed to teardown test database: ", err)
	}
}

// Configures a router to use during testing.
func SetupTestRouter(handler *handler.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return router.New(handler)
}
