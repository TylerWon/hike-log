package store_test

import (
	"errors"
	"testing"

	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type storeTestSuite struct {
	suite.Suite
	suiteDB *testutils.TestSuiteDB
	store   store.Store
}

func (suite *storeTestSuite) SetupSuite() {
	suite.suiteDB = testutils.NewTestSuiteDB(suite.T())
	suite.store = testutils.NewStore(suite.T(), suite.suiteDB)
}

func (suite *storeTestSuite) TearDownSuite() {
	testutils.TeardownStore(suite.T(), suite.store)
	suite.suiteDB.Teardown(suite.T())
}

func (suite *storeTestSuite) SetupTest() {
	suite.suiteDB.Reset(suite.T())
}

func (suite *storeTestSuite) TestCreateHike_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		CreateResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hike := testutils.CreateHikes(suite.T(), 1, suite.store, false, false)[0]
	err := store.CreateHike(&hike)

	suite.Error(err)
}

func (suite *storeTestSuite) TestCreateHike_CreatesHike() {
	hike := testutils.CreateHikes(suite.T(), 1, suite.store, false, false)[0]
	err := suite.store.CreateHike(&hike)

	suite.NoError(err)

	result, err := suite.store.GetHikeByID(hike.ID)
	suite.NoError(err)
	suite.Equal(hike, *result)
}

func (suite *storeTestSuite) TestCreateHikes_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		CreateResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hikes := testutils.CreateHikes(suite.T(), 2, suite.store, false, false)
	err := store.CreateHikes(hikes)

	suite.Error(err)
}

func (suite *storeTestSuite) TestCreateHikes_CreatesHikes() {
	hikes := testutils.CreateHikes(suite.T(), 2, suite.store, false, false)
	err := suite.store.CreateHikes(hikes)

	suite.NoError(err)

	result, err := suite.store.GetHikeByID(hikes[0].ID)
	suite.NoError(err)
	suite.Equal(hikes[0], *result)

	result, err = suite.store.GetHikeByID(hikes[1].ID)
	suite.NoError(err)
	suite.Equal(hikes[1], *result)
}

func (suite *storeTestSuite) TestGetHikeByID_ReturnsErrorWhenHikeDoesNotExist() {
	_, err := suite.store.GetHikeByID(1)
	suite.Error(err)
}

func (suite *storeTestSuite) TestGetHikeByID_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		FirstResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hike := testutils.CreateHikes(suite.T(), 1, suite.store, false, true)[0]
	_, err := store.GetHikeByID(hike.ID)

	suite.Error(err)
}

func (suite *storeTestSuite) TestGetHikeByID_ReturnsHike() {
	hike := testutils.CreateHikes(suite.T(), 1, suite.store, false, true)[0]
	result, err := suite.store.GetHikeByID(hike.ID)

	suite.NoError(err)
	suite.Equal(hike, *result)
}

func (suite *storeTestSuite) TestListHikes_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		FindResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	testutils.CreateHikes(suite.T(), 2, suite.store, true, true)
	_, err := store.ListHikes()

	suite.Error(err)
}

func (suite *storeTestSuite) TestListHikes_ReturnsHikes() {
	hikes := testutils.CreateHikes(suite.T(), 2, suite.store, true, true)
	result, err := suite.store.ListHikes()

	suite.NoError(err)
	suite.Equal(hikes[1], result[0])
	suite.Equal(hikes[0], result[1])
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(storeTestSuite))
}
