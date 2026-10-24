package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/gin-gonic/gin"
)

/*
Middleware that checks if a request is by an authenticated user.

Terminates middleware chain if request is unauthenticated or check fails.
*/
func RequireAuth(c *gin.Context) {
	authenticated, err := auth.IsAuthenticated(c)
	if !authenticated {
		// Log error if authentication fails for non-auth reason
		if err != nil {
			_, ok := errors.AsType[auth.AuthError](err)
			if !ok {
				log.Println(err)
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}

	c.Next()
}
