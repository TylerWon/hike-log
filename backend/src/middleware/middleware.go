package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/TylerWon/hike-log/backend/src/auth"
	"github.com/gin-gonic/gin"
)

const (
	RateLimit  = 10
	RateWindow = 15 * time.Minute
)

// Middleware provides reusable middleware for API routes.
type Middleware struct {
	auth        auth.Auth
	rateLimiter *RateLimiter
}

// Creates a new Middleware.
func New(auth auth.Auth) *Middleware {
	return &Middleware{
		auth:        auth,
		rateLimiter: NewRateLimiter(RateLimit, RateWindow),
	}
}

/*
Middleware that checks if a request is by an authenticated user.

Terminates middleware chain with 401 when request is unauthenticated.
*/
func (m *Middleware) RequireAuth(c *gin.Context) {
	ok, _ := m.auth.IsAuthenticated(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}
	c.Next()
}

/*
Middleware that rate limits requests by client IP.

Terminates middleware chain with 429 when the limit is exceeded.
*/
func (m *Middleware) RateLimit(c *gin.Context) {
	ip := c.ClientIP()
	if !m.rateLimiter.IsAllowed(ip) {
		slog.Warn("Rate limit exceeded", "ip", ip)
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": http.StatusText(http.StatusTooManyRequests)})
		return
	}
	c.Next()
}
