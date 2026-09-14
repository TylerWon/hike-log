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

func (suite *storeTestSuite) TestCreateModel_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		CreateResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, false)[0]
	err := store.CreateModel(&hike)

	suite.Error(err)
}

func (suite *storeTestSuite) TestCreateModel_CreatesModel() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, false)[0]
	err := suite.store.CreateModel(&hike)

	suite.NoError(err)

	result, err := suite.store.GetHikeByID(hike.ID)
	suite.NoError(err)
	suite.Equal(hike, *result)
}

func (suite *storeTestSuite) TestCreateModel_CreatesModels() {
	hikes := testutils.ConstructHikes(suite.T(), 2, suite.store, false, false)
	err := suite.store.CreateModel(hikes)

	suite.NoError(err)

	result, err := suite.store.GetHikeByID(hikes[0].ID)
	suite.NoError(err)
	suite.Equal(hikes[0], *result)

	result, err = suite.store.GetHikeByID(hikes[1].ID)
	suite.NoError(err)
	suite.Equal(hikes[1], *result)
}

func (suite *storeTestSuite) TestDeleteModel_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		DeleteResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]
	err := store.DeleteModel(&hike)
	suite.Error(err)
}

func (suite *storeTestSuite) TestDeleteModel_NoErrorWhenHikeDoesNotExist() {
	err := suite.store.DeleteModel(&models.Hike{ID: 1})
	suite.NoError(err)
}

func (suite *storeTestSuite) TestDeleteModel_DeletesModel() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]
	err := suite.store.DeleteModel(&hike)
	suite.NoError(err)

	_, err = suite.store.GetHikeByID(hike.ID)
	suite.Error(err)
}

func (suite *storeTestSuite) TestDeleteModel_DeletesModels() {
	hikes := testutils.ConstructHikes(suite.T(), 2, suite.store, false, true)
	err := suite.store.DeleteModel(&hikes)
	suite.NoError(err)

	_, err = suite.store.GetHikeByID(hikes[0].ID)
	suite.Error(err)
	_, err = suite.store.GetHikeByID(hikes[1].ID)
	suite.Error(err)
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

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]
	_, err := store.GetHikeByID(hike.ID)

	suite.Error(err)
}

func (suite *storeTestSuite) TestGetHikeByID_ReturnsHike() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]
	result, err := suite.store.GetHikeByID(hike.ID)

	suite.NoError(err)
	suite.Equal(hike, *result)
}

func (suite *storeTestSuite) TestGetPhotoByID_ReturnsErrorWhenPhotoDoesNotExist() {
	_, err := suite.store.GetPhotoByID(1)
	suite.Error(err)
}

func (suite *storeTestSuite) TestGetPhotoByID_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		FirstResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]
	photo := hike.Photos[0]
	_, err := store.GetPhotoByID(photo.ID)

	suite.Error(err)
}

func (suite *storeTestSuite) TestGetPhotoByID_ReturnsPhoto() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]
	photo := hike.Photos[0]
	result, err := suite.store.GetPhotoByID(photo.ID)

	suite.NoError(err)
	suite.Equal(hike.Photos[0], *result)
}

func (suite *storeTestSuite) TestListHikes_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		FindResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	testutils.ConstructHikes(suite.T(), 2, suite.store, true, true)
	_, err := store.ListHikes()

	suite.Error(err)
}

func (suite *storeTestSuite) TestListHikes_ReturnsHikes() {
	hikes := testutils.ConstructHikes(suite.T(), 2, suite.store, true, true)
	result, err := suite.store.ListHikes()

	suite.NoError(err)
	suite.Equal(hikes[1], result[0])
	suite.Equal(hikes[0], result[1])
}

func (suite *storeTestSuite) TestUpdateModel_ReturnsErrorWhenDBErrors() {
	mockDatabase := database.MockDatabase{
		UpdatesResult: &gorm.DB{Error: errors.New("Something went wrong")},
	}
	store := store.NewTestStore(&mockDatabase)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]

	hike.Difficulty = 5
	hike.TrailName = "Other trail name"
	err := store.UpdateModel(&hike)
	suite.Error(err)
}

func (suite *storeTestSuite) TestUpdateModel_NoErrorWhenHikeDoesNotExist() {
	err := suite.store.UpdateModel(&models.Hike{ID: 1})
	suite.NoError(err)
}

func (suite *storeTestSuite) TestUpdateModel_UpdatesModel() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]

	hike.Difficulty = 5
	hike.TrailName = "Other trail name"
	err := suite.store.UpdateModel(&hike)
	suite.NoError(err)

	updated, err := suite.store.GetHikeByID(hike.ID)
	suite.NoError(err)
	suite.Equal(hike, *updated)
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(storeTestSuite))
}
