package auth

import "github.com/gin-gonic/gin"

type MockAuth struct {
	CreateSessionCookieValueResult    string
	GenerateHmacSha256SignatureResult string
	IsAuthenticatedResult             bool
	IsAuthenticatedError              error
	VerifyPasswordResult              bool
}

func (m *MockAuth) CreateSessionCookieValue() string {
	return m.CreateSessionCookieValueResult
}

func (m *MockAuth) GenerateHmacSha256Signature(payload string) string {
	return m.GenerateHmacSha256SignatureResult
}

func (m *MockAuth) IsAuthenticated(c *gin.Context) (bool, error) {
	return m.IsAuthenticatedResult, m.IsAuthenticatedError
}

func (m *MockAuth) SetSessionCookie(c *gin.Context) {}

func (m *MockAuth) VerifyPassword(password string) bool {
	return m.VerifyPasswordResult
}
