package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

/*
Creates a new Gin context. Optionally sets a cookie on the request if provided.
*/
func NewGinContext(t *testing.T, cookie *http.Cookie) *gin.Context {
	t.Helper()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if cookie != nil {
		req.AddCookie(cookie)
	}
	c.Request = req

	return c
}
