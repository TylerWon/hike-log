package middleware

import (
	"net/http"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/gin-gonic/gin"
)

// Middleware provides reusable middleware for API routes.
type Middleware struct {
	auth auth.Auth
}

// Creates a new Middleware.
func New(auth auth.Auth) *Middleware {
	return &Middleware{auth}
}

/*
Middleware that checks if a request is by an authenticated user.

Terminates middleware chain if request is unauthenticated.
*/
func (m *Middleware) RequireAuth(c *gin.Context) {
	ok, _ := m.auth.IsAuthenticated(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}
	c.Next()
}
