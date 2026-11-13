package testutils

import (
	"net/http"
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/config"
)

// Creates a session cookie that can be attached to requests.
func CreateSessionCookie(t *testing.T, authService auth.Auth) *http.Cookie {
	t.Helper()

	return &http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: authService.CreateSessionCookieValue(),
	}
}

// Creates an Auth to use during testing. Arguments are optional and can be provided to override default values.
func NewAuth(t *testing.T, authConfig config.AuthConfig) auth.Auth {
	t.Helper()

	testConfig := NewTestConfig()

	if authConfig.AdminPasswordHash != "" {
		testConfig.Auth.AdminPasswordHash = authConfig.AdminPasswordHash
	}
	if authConfig.SessionSecret != "" {
		testConfig.Auth.SessionSecret = authConfig.SessionSecret
	}

	auth, err := auth.New(testConfig.Auth)
	if err != nil {
		t.Fatal("Failed to setup auth: ", err)
	}

	return auth
}
