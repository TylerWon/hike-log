package database_test

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/stretchr/testify/suite"
)

type databaseTestSuite struct {
	suite.Suite
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(databaseTestSuite))
}

func (suite *databaseTestSuite) TestConnect_ReturnsErrorWhenConnectionFails() {
	dbConfig := testutils.NewTestConfig().Database
	dbConfig.DBHost = "127.0.0.1"
	dbConfig.DBPort = "1"

	db, err := database.Connect(dbConfig)
	suite.Require().Error(err)
	suite.Require().Nil(db)
}

func (suite *databaseTestSuite) TestConnect_CreatesDatabase() {
	dbConfig := testutils.NewTestConfig().Database
	dbConfig.DBName = "postgres" // postgres database guaranteed to exist if database is running; same can't be said for other database

	db, err := database.Connect(dbConfig)
	suite.Require().NoError(err)
	suite.Require().NotNil(db)

	sqlDB, err := db.DB()
	suite.Require().NoError(err)
	suite.Require().NoError(sqlDB.Close())
}
