package database_test

import (
	"context"
	"errors"
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

	suite.Require().NoError(db.Close())
}
func (suite *databaseTestSuite) TestClose_ReturnsErrorWhenDatabaseErrors() {
	db := &database.MockDatabase{CloseError: errors.New("Something went wrong")}
	err := db.Close()
	suite.Require().Error(err)
}

func (suite *databaseTestSuite) TestClose_ClosesDatabaseConnection() {
	dbConfig := testutils.NewTestConfig().Database
	dbConfig.DBName = "postgres"

	db, err := database.Connect(dbConfig)
	suite.Require().NoError(err)
	suite.Require().NotNil(db)

	suite.Require().NoError(db.Close())
	suite.Require().Error(db.Ping(context.Background()))
}

func (suite *databaseTestSuite) TestPing_ReturnsErrorWhenDatabaseErrors() {
	db := &database.MockDatabase{PingError: errors.New("Something went wrong")}
	err := db.Ping(context.Background())
	suite.Require().Error(err)
}

func (suite *databaseTestSuite) TestPing_ReturnsNil() {
	dbConfig := testutils.NewTestConfig().Database
	dbConfig.DBName = "postgres"

	db, err := database.Connect(dbConfig)
	suite.Require().NoError(err)
	suite.Require().NotNil(db)

	suite.Require().NoError(db.Ping(context.Background()))
	suite.Require().NoError(db.Close())
}
