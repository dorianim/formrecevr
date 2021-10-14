package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewApiTest(groupPath string) (app *gin.Engine, router *gin.RouterGroup) {
	gin.SetMode(gin.TestMode)
	app = gin.New()
	router = app.Group(groupPath)
	return app, router
}

// Performs API request with empty request body.
// See https://medium.com/@craigchilds94/testing-gin-json-responses-1f258ce3b0b1
func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// Performs API request including request body as string.
func PerformRequestWithBody(r http.Handler, method, path, contentType string, body string) *httptest.ResponseRecorder {
	reader := strings.NewReader(body)
	req, _ := http.NewRequest(method, path, reader)
	req.Header.Add("Content-Type", contentType)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
