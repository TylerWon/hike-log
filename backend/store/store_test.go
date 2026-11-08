package store_test

import (
	"errors"
	"testing"

	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/models"
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

func (suite *storeTestSuite) TestCreateRecord_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		CreateResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, false)[0]
	err := store.CreateRecord(&hike)

	suite.Require().Error(err)
}

func (suite *storeTestSuite) TestCreateRecord_CreatesModel() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, false)[0]
	err := suite.store.CreateRecord(&hike)

	suite.Require().NoError(err)

	result, err := suite.store.GetHikeByID(hike.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(hike, *result)
}

func (suite *storeTestSuite) TestCreateRecord_CreatesModels() {
	hikes := testutils.ConstructHikes(suite.T(), 2, suite.store, false, false)
	err := suite.store.CreateRecord(hikes)

	suite.Require().NoError(err)

	result, err := suite.store.GetHikeByID(hikes[0].ID)
	suite.Require().NoError(err)
	suite.Require().Equal(hikes[0], *result)

	result, err = suite.store.GetHikeByID(hikes[1].ID)
	suite.Require().NoError(err)
	suite.Require().Equal(hikes[1], *result)
}

func (suite *storeTestSuite) TestDeleteRecord_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		DeleteResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]
	err := store.DeleteRecord(&hike)
	suite.Require().Error(err)
}

func (suite *storeTestSuite) TestDeleteRecord_NoErrorWhenHikeDoesNotExist() {
	err := suite.store.DeleteRecord(&models.Hike{ID: 1})
	suite.Require().NoError(err)
}

func (suite *storeTestSuite) TestDeleteRecord_DeletesModel() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]
	err := suite.store.DeleteRecord(&hike)
	suite.Require().NoError(err)

	_, err = suite.store.GetHikeByID(hike.ID)
	suite.Require().Error(err)
}

func (suite *storeTestSuite) TestDeleteRecord_DeletesModels() {
	hikes := testutils.ConstructHikes(suite.T(), 2, suite.store, false, true)
	err := suite.store.DeleteRecord(&hikes)
	suite.Require().NoError(err)

	_, err = suite.store.GetHikeByID(hikes[0].ID)
	suite.Require().Error(err)
	_, err = suite.store.GetHikeByID(hikes[1].ID)
	suite.Require().Error(err)
}

func (suite *storeTestSuite) TestGetHikeByID_ReturnsErrorWhenHikeDoesNotExist() {
	_, err := suite.store.GetHikeByID(1)
	suite.Require().Error(err)
}

func (suite *storeTestSuite) TestGetHikeByID_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		FirstResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]
	_, err := store.GetHikeByID(hike.ID)

	suite.Require().Error(err)
}

func (suite *storeTestSuite) TestGetHikeByID_ReturnsHike() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]
	result, err := suite.store.GetHikeByID(hike.ID)

	suite.Require().NoError(err)
	suite.Require().Equal(hike, *result)
}

func (suite *storeTestSuite) TestGetPhotoByID_ReturnsErrorWhenPhotoDoesNotExist() {
	_, err := suite.store.GetPhotoByID(1)
	suite.Require().Error(err)
}

func (suite *storeTestSuite) TestGetPhotoByID_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		FirstResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]
	photo := hike.Photos[0]
	_, err := store.GetPhotoByID(photo.ID)

	suite.Require().Error(err)
}

func (suite *storeTestSuite) TestGetPhotoByID_ReturnsPhoto() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]
	photo := hike.Photos[0]
	result, err := suite.store.GetPhotoByID(photo.ID)

	suite.Require().NoError(err)
	suite.Require().Equal(hike.Photos[0], *result)
}

func (suite *storeTestSuite) TestListHikes_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		FindResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	testutils.ConstructHikes(suite.T(), 2, suite.store, true, true)
	_, err := store.ListHikes()

	suite.Require().Error(err)
}

func (suite *storeTestSuite) TestListHikes_ReturnsHikes() {
	hikes := testutils.ConstructHikes(suite.T(), 2, suite.store, true, true)
	result, err := suite.store.ListHikes()

	suite.Require().NoError(err)
	suite.Require().Equal(hikes[1], result[0])
	suite.Require().Equal(hikes[0], result[1])
}

func (suite *storeTestSuite) TestUpdateRecord_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		UpdatesResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]

	hike.Difficulty = 5
	hike.TrailName = "Other trail name"
	err := store.UpdateRecord(&hike)
	suite.Require().Error(err)
}

func (suite *storeTestSuite) TestUpdateRecord_NoErrorWhenHikeDoesNotExist() {
	err := suite.store.UpdateRecord(&models.Hike{ID: 1})
	suite.Require().NoError(err)
}

func (suite *storeTestSuite) TestUpdateRecord_UpdatesModel() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]

	hike.Difficulty = 5
	hike.TrailName = "Other trail name"
	err := suite.store.UpdateRecord(&hike)
	suite.Require().NoError(err)

	updated, err := suite.store.GetHikeByID(hike.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(hike, *updated)
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(storeTestSuite))
}
