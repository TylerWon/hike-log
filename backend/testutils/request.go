package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Serializes a JSON request body.
func SerializeJSONRequestBody(t *testing.T, body map[string]any) *bytes.Reader {
	rawBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(rawBody)
}

// Sends a request to the provided endpoint and returns the response.
func SendRequest(router *gin.Engine, method string, endpoint string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, endpoint, body)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}
