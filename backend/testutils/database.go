package testutils

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/TylerWon/hike-log/backend/database"
)

// TestSuiteDB holds a shared database for a test suite.
type TestSuiteDB struct {
	DBName  string
	adminDB database.Database
	testDB  database.Database
}

// Creates a dedicated database for a test suite.
func NewTestSuiteDB(t *testing.T) *TestSuiteDB {
	t.Helper()

	adminDB := connectToAdminDB(t)
	dbName := strings.ToLower(strings.ReplaceAll(t.Name(), "/", "_"))

	var exists int
	result := adminDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", dbName).Scan(&exists)
	if result.Error != nil {
		t.Fatalf("Unable to verify if old test suite database (%s) was deleted: %v", dbName, result.Error)
	} else if exists != 0 {
		result := adminDB.Exec(fmt.Sprintf("DROP DATABASE %s", dbName))
		if result.Error != nil {
			t.Fatalf("Failed to delete old test suite database (%s): %v", dbName, result.Error)
		}
	}

	result = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if result.Error != nil {
		t.Fatalf("Failed to create test suite database (%s): %v", dbName, result.Error)
	}

	dbConfig := createDBConfig(dbName)
	testDB, err := database.Connect(dbConfig)
	if err != nil {
		t.Fatalf("Failed to connect to test suite database (%s): %v", dbName, err)
	}

	return &TestSuiteDB{
		adminDB: adminDB,
		testDB:  testDB,
		DBName:  dbName,
	}
}

// Truncates all application tables from the test suite database.
func (suiteDB *TestSuiteDB) Reset(t *testing.T) {
	t.Helper()

	result := suiteDB.testDB.Exec("TRUNCATE TABLE photos, hikes RESTART IDENTITY CASCADE")
	if result.Error != nil {
		t.Fatalf("Failed to reset %s database: %v", suiteDB.DBName, result.Error)
	}
}

// Drops the suite database and closes all connections.
func (suiteDB *TestSuiteDB) Teardown(t *testing.T) {
	t.Helper()

	closeDB(t, suiteDB.testDB)

	var exists int
	suiteDB.adminDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", suiteDB.DBName).Scan(&exists)
	if exists != 0 {
		result := suiteDB.adminDB.Exec(fmt.Sprintf("DROP DATABASE %s", suiteDB.DBName))
		if result.Error != nil {
			t.Fatalf("Failed to delete test suite database (%s): %v", suiteDB.DBName, result.Error)
		}
	}

	closeDB(t, suiteDB.adminDB)
}

func connectToAdminDB(t *testing.T) database.Database {
	t.Helper()

	db, err := database.Connect(createDBConfig("postgres"))
	if err != nil {
		t.Fatal("Failed to connect to 'postgres' database: ", err)
	}

	return db
}

func createDBConfig(dbName string) database.DBConfig {
	return database.DBConfig{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbName:     dbName,
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
	}
}

func closeDB(t *testing.T, db database.Database) {
	t.Helper()

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal("Failed to get database connection: ", err)
	}

	if err := sqlDB.Close(); err != nil {
		t.Fatal("Failed to close database connection: ", err)
	}
}
