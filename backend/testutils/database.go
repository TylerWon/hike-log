package testutils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/database"
)

// TestSuiteDB holds a shared database for a test suite.
type TestSuiteDB struct {
	TestDB       database.Database
	adminDB      database.Database
	testDBConfig config.DBConfig
}

// Creates a dedicated database for a test suite.
func NewTestSuiteDB(t *testing.T) *TestSuiteDB {
	t.Helper()

	testConfig := NewTestConfig()

	adminDBConfig := testConfig.Database
	adminDBConfig.DBName = "postgres"
	adminDB, err := database.Connect(adminDBConfig)
	if err != nil {
		t.Fatal("Failed to connect to 'postgres' database: ", err)
	}

	testDBConfig := testConfig.Database
	testSuiteName := strings.ToLower(strings.ReplaceAll(t.Name(), "/", "_"))
	testDBConfig.DBName = testDBConfig.DBName + "__" + testSuiteName

	var exists int
	result := adminDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", testDBConfig.DBName).Scan(&exists)
	if result.Error != nil {
		t.Fatalf("Unable to verify if old test suite database (%s) was deleted: %v", testDBConfig.DBName, result.Error)
	} else if exists != 0 {
		result := adminDB.Exec(fmt.Sprintf("DROP DATABASE %s", testDBConfig.DBName))
		if result.Error != nil {
			t.Fatalf("Failed to delete old test suite database (%s): %v", testDBConfig.DBName, result.Error)
		}
	}

	result = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBConfig.DBName))
	if result.Error != nil {
		t.Fatalf("Failed to create test suite database (%s): %v", testDBConfig.DBName, result.Error)
	}

	testDB, err := database.Connect(testDBConfig)
	if err != nil {
		t.Fatalf("Failed to connect to test suite database (%s): %v", testDBConfig.DBName, err)
	}

	return &TestSuiteDB{
		adminDB:      adminDB,
		TestDB:       testDB,
		testDBConfig: testDBConfig,
	}
}

// Truncates all application tables from the test suite database.
func (suiteDB *TestSuiteDB) Reset(t *testing.T) {
	t.Helper()

	result := suiteDB.TestDB.Exec("TRUNCATE TABLE photos, hikes RESTART IDENTITY CASCADE")
	if result.Error != nil {
		t.Fatalf("Failed to reset %s database: %v", suiteDB.testDBConfig.DBName, result.Error)
	}
}

// Drops the suite database and closes all connections.
func (suiteDB *TestSuiteDB) Teardown(t *testing.T) {
	t.Helper()

	err := suiteDB.TestDB.Close()
	if err != nil {
		t.Fatal("Failed to close test database connection: ", err)
	}

	var exists int
	suiteDB.adminDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", suiteDB.testDBConfig.DBName).Scan(&exists)
	if exists != 0 {
		result := suiteDB.adminDB.Exec(fmt.Sprintf("DROP DATABASE %s", suiteDB.testDBConfig.DBName))
		if result.Error != nil {
			t.Fatalf("Failed to delete test suite database (%s): %v", suiteDB.testDBConfig.DBName, result.Error)
		}
	}

	err = suiteDB.adminDB.Close()
	if err != nil {
		t.Fatal("Failed to close admin database connection: ", err)
	}
}
