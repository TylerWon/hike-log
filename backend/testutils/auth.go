package testutils

import (
	"net/http"
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
)

// Creates a session cookie that can be attached to requests.
func CreateSessionCookie(authService auth.Auth) *http.Cookie {
	return &http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: authService.CreateSessionCookieValue(),
	}
}

// Creates an Auth to use during testing.
func NewAuth(t *testing.T, adminPasswordHash string, sessionSecret string) auth.Auth {
	t.Helper()

	auth, err := auth.New(adminPasswordHash, sessionSecret)
	if err != nil {
		t.Fatal("Failed to setup auth: ", err)
	}

	return auth
}
