package testutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/TylerWon/todo-app/backend/database"
	"github.com/TylerWon/todo-app/backend/handler"
	"github.com/TylerWon/todo-app/backend/router"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
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

	// Setup the test database
	dbConfig.DbName = testDbName
	db, err = database.Setup(dbConfig)
	if err != nil {
		t.Fatal("Failed to setup test database: ", err)
	}

	return db
}

func TeardownTestDB(t *testing.T, db *gorm.DB) {
	sqlDb, err := db.DB()
	if err != nil {
		t.Fatal("Failed to teardown test database: ", err)
	}

	sqlDb.Close()
}

func SetupTestRouter(handler *handler.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return router.New(handler)
}
