package handler_test

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/s3"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type handlerTestSuite struct {
	suite.Suite
	auth       auth.Auth
	middleware *middleware.Middleware
	s3Client   s3.S3Client
	store      store.Store
	suiteDB    *testutils.TestSuiteDB
	router     *gin.Engine
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
	suite.auth = testutils.NewAuth(suite.T(), testutils.ADMIN_PASSWORD_HASH, testutils.SESSION_SECRET)
	suite.middleware = middleware.New(suite.auth)
	handler := handler.New(suite.store, suite.s3Client, suite.auth)
	suite.router = testutils.NewRouter(handler, suite.middleware)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
