package testutils

import (
	"net/http"
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
)

const (
	// bcrypt hash for admin password
	ADMIN_PASSWORD_HASH = "$2a$10$haTJ.11uqjdIRUD4yxMSt.HfJiOL9.WjyjH1D2C.qv8okP0c0wuXC"
	// Raw password string
	ADMIN_PASSWORD = "admin"
	// Random string for generating session cookies
	SESSION_SECRET = "vDuTmQE8HJj_SsRT32y8Wz8yH97FoRCT5VUiNJlfsGI"
)

// Creates a session cookie that can be attached to requests.
func CreateSessionCookie(authService auth.Auth) *http.Cookie {
	return &http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: authService.CreateSessionCookieValue(),
	}
}

// Creates an Auth to use during testing. Arguments are optional and can be provided to override default values.
func NewAuth(t *testing.T, adminPasswordHash *string, sessionSecret *string) auth.Auth {
	t.Helper()

	hash := ADMIN_PASSWORD_HASH
	secret := SESSION_SECRET

	if adminPasswordHash != nil {
		hash = *adminPasswordHash
	}

	if sessionSecret != nil {
		secret = *sessionSecret
	}

	auth, err := auth.New(hash, secret)
	if err != nil {
		t.Fatal("Failed to setup auth: ", err)
	}

	return auth
}
