package testutils

import (
	"net/http"
	"os"

	"github.com/TylerWon/hike-log/backend/auth"
)

func CreateValidSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: auth.CreateSessionCookieValue(os.Getenv("SESSION_SECRET")),
	}
}
