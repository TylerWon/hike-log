package handler_test

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/aws"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type handlerTestSuite struct {
	suite.Suite
	router   *gin.Engine
	suiteDB  *testutils.TestSuiteDB
	store    store.Store
	s3Client aws.S3Client
}

func (suite *handlerTestSuite) SetupSuite() {
	suite.suiteDB = testutils.NewTestSuiteDB(suite.T())
	suite.store = testutils.NewStore(suite.T(), suite.suiteDB)
}

func (suite *handlerTestSuite) TearDownSuite() {
	testutils.TeardownStore(suite.T(), suite.store)
	suite.suiteDB.Teardown(suite.T())
}

func (suite *handlerTestSuite) SetupTest() {
	suite.suiteDB.Reset(suite.T())
	suite.s3Client = testutils.NewS3Client(suite.T())
	handler := handler.New(suite.store, suite.s3Client)
	suite.router = testutils.NewRouter(suite.T(), handler)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
