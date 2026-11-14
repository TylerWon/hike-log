package handler_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/TylerWon/hike-log/backend/database"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/stretchr/testify/suite"
)

type healthHandlerTestSuite struct {
	suite.Suite
}

func TestHealthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(healthHandlerTestSuite))
}

func (suite *healthHandlerTestSuite) TestNewHealthHandler_CreatesHealthHandler() {
	h := handler.NewHealthHandler(&database.MockDatabase{})
	suite.Require().NotNil(h)
}

func (suite *healthHandlerTestSuite) TestHealthCheck_Returns200() {
	h := testutils.NewHealthHandler(suite.T(), &database.MockDatabase{})
	router := testutils.NewHealthRouter(suite.T(), h)

	res := testutils.SendRequest(suite.T(), router, http.MethodGet, "/health", nil)

	suite.Require().Equal(http.StatusOK, res.Code)
}

func (suite *healthHandlerTestSuite) TestReadyCheck_Returns503WhenDatabaseErrors() {
	h := testutils.NewHealthHandler(suite.T(), &database.MockDatabase{
		PingError: errors.New("Something went wrong"),
	})
	router := testutils.NewHealthRouter(suite.T(), h)

	res := testutils.SendRequest(suite.T(), router, http.MethodGet, "/ready", nil)

	suite.Require().Equal(http.StatusServiceUnavailable, res.Code)
}

func (suite *healthHandlerTestSuite) TestReadyCheck_Returns200WhenPingSucceeds() {
	h := testutils.NewHealthHandler(suite.T(), &database.MockDatabase{})
	router := testutils.NewHealthRouter(suite.T(), h)

	res := testutils.SendRequest(suite.T(), router, http.MethodGet, "/ready", nil)

	suite.Require().Equal(http.StatusOK, res.Code)
}
