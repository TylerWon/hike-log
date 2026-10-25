package middleware_test

import (
	"net/http"
	"testing"

	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/stretchr/testify/suite"
)

type authTestSuite struct {
	suite.Suite
}

func (suite *authTestSuite) TestRequireAuth_Returns500ErrorWhenSessionSecretIsMissing() {
	c := testutils.NewGinContext(testutils.CreateValidSessionCookie())

	suite.T().Setenv("SESSION_SECRET", "")

	middleware.RequireAuth(c)

	res := c.Writer
	suite.Equal(http.StatusInternalServerError, res.Status())
	suite.True(c.IsAborted())
}

func (suite *authTestSuite) TestRequireAuth_Returns401ErrorAndAbortsWhenUnauthenticated() {
	c := testutils.NewGinContext(nil) // Missing session cookie

	middleware.RequireAuth(c)

	res := c.Writer
	suite.Equal(http.StatusUnauthorized, res.Status())
	suite.True(c.IsAborted())
}

func (suite *authTestSuite) TestRequireAuth_PassesWhenAuthenticated() {
	c := testutils.NewGinContext(testutils.CreateValidSessionCookie())
	middleware.RequireAuth(c)
	suite.False(c.IsAborted())
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(authTestSuite))
}
