package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dorianim/formrecevr/internal/api"
	"github.com/dorianim/formrecevr/internal/config"
)

func TestPostForm(t *testing.T) {

	genericConfig := &config.Config{
		Listen: config.ListenConfig{
			Host: "127.0.0.1",
			Port: 8088,
		},
		Forms: map[string]*config.FormConfig{
			"example": {
				Enabled: true,
				Targets: []*config.TargetConfig{
					{
						Enabled:     true,
						Template:    "../../templates/default.html",
						ShoutrrrURL: "generic://127.0.0.1", // send webhook to the local testserver
					},
					{
						Enabled: false,
					},
				},
			},
		},
	}

	t.Run("successful request urlencoded", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request parameters
			assert.Equal(t, req.URL.String(), "/")
			// Send response to be tested
			rw.WriteHeader(200)
		}))
		defer server.Close()

		shoutrrrURL := fmt.Sprintf("generic://%s?disabletls=yes", strings.Replace(server.URL, "http://", "", -1))
		genericConfig.Forms["example"].Targets[0].ShoutrrrURL = shoutrrrURL

		config.SetConfig(genericConfig)

		app, router := NewApiTest("/f")
		api.PostForm(router)
		r := PerformRequestWithBody(app, "POST", "/f/example", "application/x-www-form-urlencoded", "a=a&b=b")
		assert.Equal(t, http.StatusOK, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Success", response.Message)
	})

	t.Run("successful request multipart", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request parameters
			assert.Equal(t, req.URL.String(), "/")
			// Send response to be tested
			rw.WriteHeader(200)
		}))
		defer server.Close()

		shoutrrrURL := fmt.Sprintf("generic://%s?disabletls=yes", strings.Replace(server.URL, "http://", "", -1))
		genericConfig.Forms["example"].Targets[0].ShoutrrrURL = shoutrrrURL

		config.SetConfig(genericConfig)

		app, router := NewApiTest("/f")
		api.PostForm(router)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormField("a")
		io.Copy(part, strings.NewReader("a"))
		writer.Close()

		r := PerformRequestWithBody(app, "POST", "/f/example", writer.FormDataContentType(), body.String())
		assert.Equal(t, http.StatusOK, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Success", response.Message)
	})

	t.Run("form not found", func(t *testing.T) {
		config.SetConfig(genericConfig)

		app, router := NewApiTest("/f")
		api.PostForm(router)
		r := PerformRequest(app, "POST", "/f/notExistent")
		assert.Equal(t, http.StatusNotFound, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Form not found", response.Message)
	})

	t.Run("unsupported content type", func(t *testing.T) {
		config.SetConfig(genericConfig)

		app, router := NewApiTest("/f")
		api.PostForm(router)
		r := PerformRequestWithBody(app, "POST", "/f/example", "unsupported", "a=a&b=b")
		assert.Equal(t, http.StatusBadRequest, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Unsupported Content-Type", response.Message)
	})

	t.Run("malformed form", func(t *testing.T) {
		config.SetConfig(genericConfig)

		app, router := NewApiTest("/f")
		api.PostForm(router)
		r := PerformRequestWithBody(app, "POST", "/f/example", "multipart/form-data", "")
		assert.Equal(t, http.StatusBadRequest, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Malformed form data", response.Message)
	})

	t.Run("form not enabled", func(t *testing.T) {
		config.SetConfig(&config.Config{
			Listen: config.ListenConfig{
				Host: "0.0.0.0",
				Port: 8088,
			},
			Forms: map[string]*config.FormConfig{
				"example": {
					Enabled: false,
				},
			},
		})

		app, router := NewApiTest("/f")
		api.PostForm(router)
		r := PerformRequestWithBody(app, "POST", "/f/example", "application/x-www-form-urlencoded", "a=a&b=b")
		assert.Equal(t, http.StatusNotFound, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Form not found", response.Message)
	})

	t.Run("errors in all targets", func(t *testing.T) {
		config.SetConfig(&config.Config{
			Listen: config.ListenConfig{
				Host: "0.0.0.0",
				Port: 8088,
			},
			Forms: map[string]*config.FormConfig{
				"example": {
					Enabled: true,
					Targets: []*config.TargetConfig{
						{
							Enabled:     true,
							Template:    "../../templates/default.html",
							ShoutrrrURL: "malformed://target",
						},
					},
				},
			},
		})

		app, router := NewApiTest("/f")
		api.PostForm(router)
		r := PerformRequestWithBody(app, "POST", "/f/example", "application/x-www-form-urlencoded", "a=a&b=b")
		assert.Equal(t, http.StatusInternalServerError, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Internal server error", response.Message)
	})

	t.Run("error in all templates", func(t *testing.T) {
		config.SetConfig(&config.Config{
			Listen: config.ListenConfig{
				Host: "0.0.0.0",
				Port: 8088,
			},
			Forms: map[string]*config.FormConfig{
				"example": {
					Enabled: true,
					Targets: []*config.TargetConfig{
						{
							Enabled:     true,
							Template:    "../../testdata/invalid-template.html",
							ShoutrrrURL: "generic://example.com", // send webhook to example.com,
						},
					},
				},
			},
		})

		app, router := NewApiTest("/f")
		api.PostForm(router)
		r := PerformRequestWithBody(app, "POST", "/f/example", "application/x-www-form-urlencoded", "a=a&b=b")
		assert.Equal(t, http.StatusInternalServerError, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Internal server error", response.Message)
	})
}
