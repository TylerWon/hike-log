package middleware_test

import (
	"net/http"
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/stretchr/testify/suite"
)

type middlewareTestSuite struct {
	suite.Suite
	auth       auth.Auth
	middleware *middleware.Middleware
}

func (suite *middlewareTestSuite) SetupSuite() {
	suite.auth = testutils.NewAuth(suite.T(), config.AuthConfig{})
	suite.middleware = middleware.New(suite.auth)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(middlewareTestSuite))
}

func (suite *middlewareTestSuite) TestNew_CreatesMiddleware() {
	m := middleware.New(suite.auth)
	suite.Require().NotNil(m)
}

func (suite *middlewareTestSuite) TestRequireAuth_Returns401ErrorAndAbortsWhenUnauthenticated() {
	c := testutils.NewGinContext(suite.T(), nil) // Missing session cookie

	suite.middleware.RequireAuth(c)

	res := c.Writer
	suite.Require().Equal(http.StatusUnauthorized, res.Status())
	suite.Require().True(c.IsAborted())
}

func (suite *middlewareTestSuite) TestRequireAuth_PassesWhenAuthenticated() {
	c := testutils.NewGinContext(suite.T(), testutils.CreateSessionCookie(suite.T(), suite.auth))

	suite.middleware.RequireAuth(c)

	suite.Require().False(c.IsAborted())
}
