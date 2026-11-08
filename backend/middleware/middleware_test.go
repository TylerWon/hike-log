package middleware_test

import (
	"net/http"
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
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
	suite.auth = testutils.NewAuth(suite.T(), nil, nil)
	suite.middleware = middleware.New(suite.auth)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(middlewareTestSuite))
}

func (suite *middlewareTestSuite) TestRequireAuth_Returns401ErrorAndAbortsWhenUnauthenticated() {
	c := testutils.NewGinContext(nil) // Missing session cookie

	suite.middleware.RequireAuth(c)

	res := c.Writer
	suite.Require().Equal(http.StatusUnauthorized, res.Status())
	suite.Require().True(c.IsAborted())
}

func (suite *middlewareTestSuite) TestRequireAuth_PassesWhenAuthenticated() {
	c := testutils.NewGinContext(testutils.CreateSessionCookie(suite.auth))

	suite.middleware.RequireAuth(c)

	suite.Require().False(c.IsAborted())
}
