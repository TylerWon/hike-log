package middleware_test

import (
	"testing"
	"time"

	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/stretchr/testify/suite"
)

type rateLimiterTestSuite struct {
	suite.Suite
}

func TestRateLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(rateLimiterTestSuite))
}

func (suite *rateLimiterTestSuite) TestIsAllowed_ReturnsTrueWhenUnderOrAtLimit() {
	limiter := middleware.NewRateLimiter(2, time.Minute)

	suite.Require().True(limiter.IsAllowed("1.1.1.1"))
	suite.Require().True(limiter.IsAllowed("1.1.1.1"))
}

func (suite *rateLimiterTestSuite) TestIsAllowed_ReturnsFalseWhenOverLimit() {
	limiter := middleware.NewRateLimiter(2, time.Minute)

	suite.Require().True(limiter.IsAllowed("1.1.1.1"))
	suite.Require().True(limiter.IsAllowed("1.1.1.1"))
	suite.Require().False(limiter.IsAllowed("1.1.1.1"))
}

func (suite *rateLimiterTestSuite) TestIsAllowed_ReturnsTrueAfterWindowExpires() {
	mockNow := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	mockGetNow := func() time.Time { return mockNow }
	limiter := middleware.NewTestRateLimiter(1, time.Minute, mockGetNow)

	suite.Require().True(limiter.IsAllowed("1.1.1.1"))
	suite.Require().False(limiter.IsAllowed("1.1.1.1"))

	mockNow = mockNow.Add(time.Minute + time.Second) // Set current time to after window expires
	suite.Require().True(limiter.IsAllowed("1.1.1.1"))
}

func (suite *rateLimiterTestSuite) TestIsAllowed_IsolatesLimitPerIP() {
	limiter := middleware.NewRateLimiter(1, time.Minute)

	suite.Require().True(limiter.IsAllowed("1.1.1.1"))
	suite.Require().False(limiter.IsAllowed("1.1.1.1"))
	suite.Require().True(limiter.IsAllowed("2.2.2.2"))
}

func (suite *rateLimiterTestSuite) TestRemoveExpired_RemovesIPsWithNoRecentRequests() {
	mockNow := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	mockGetNow := func() time.Time { return mockNow }
	limiter := middleware.NewTestRateLimiter(1, time.Minute, mockGetNow)

	suite.Require().True(limiter.IsAllowed("1.1.1.1"))
	suite.Require().True(limiter.IsAllowed("2.2.2.2"))

	mockNow = mockNow.Add(time.Minute + time.Second)
	suite.Require().True(limiter.IsAllowed("2.2.2.2")) // Keep this IP active within the new window

	limiter.RemoveExpired()

	suite.Require().False(limiter.HasAttempts("1.1.1.1"))
	suite.Require().True(limiter.HasAttempts("2.2.2.2"))
}
