package api_test

import (
	"net/http"
	"net/http/httptest"

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
