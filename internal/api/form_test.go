package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dorianim/formrecevr/internal/api"
	"github.com/dorianim/formrecevr/internal/config"
	"github.com/dorianim/formrecevr/internal/template"
)

func TestPostForm(t *testing.T) {

	tmpTemplateDir := "../../testdata/tmp"
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
						Template:    fmt.Sprintf("%s/default.html", tmpTemplateDir),
						ShoutrrrURL: "generic://127.0.0.1?disabletls=yes", // send webhook to the local testserver
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
		template.CreateDefaultTemplate(tmpTemplateDir)

		app, router := NewApiTest("/f")
		api.PostForm(router)
		r := PerformRequestWithBody(app, "POST", "/f/example", "application/x-www-form-urlencoded", "a=a&b=b")
		assert.Equal(t, http.StatusOK, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Success", response.Message)
		os.RemoveAll(tmpTemplateDir)
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
		template.CreateDefaultTemplate(tmpTemplateDir)

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
		os.RemoveAll(tmpTemplateDir)
	})

	t.Run("successful request urlencoded with shoutrrr template", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Test request parameters
			assert.Equal(t, req.URL.String(), "/")
			// Send response to be tested
			rw.WriteHeader(200)
		}))
		defer server.Close()

		tmpConfig := &config.Config{
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
							Template:    fmt.Sprintf("%s/default.html", tmpTemplateDir),
							ShoutrrrURL: "generic://{{ print .c }}?disabletls=yes", // send webhook to templatable hook
						},
						{
							Enabled: false,
						},
					},
				},
			},
		}

		config.SetConfig(tmpConfig)
		template.CreateDefaultTemplate(tmpTemplateDir)

		app, router := NewApiTest("/f")
		api.PostForm(router)
		formData := fmt.Sprintf("a=a&b=b&c=%s", strings.Replace(server.URL, "http://", "", -1))
		r := PerformRequestWithBody(app, "POST", "/f/example", "application/x-www-form-urlencoded", formData)
		assert.Equal(t, http.StatusOK, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Success", response.Message)
		os.RemoveAll(tmpTemplateDir)
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
							Template:    fmt.Sprintf("%s/default.html", tmpTemplateDir),
							ShoutrrrURL: "malformed://target",
						},
					},
				},
			},
		})
		template.CreateDefaultTemplate(tmpTemplateDir)

		app, router := NewApiTest("/f")
		api.PostForm(router)
		r := PerformRequestWithBody(app, "POST", "/f/example", "application/x-www-form-urlencoded", "a=a&b=b")
		assert.Equal(t, http.StatusInternalServerError, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Internal server error", response.Message)
		os.RemoveAll(tmpTemplateDir)
	})

	t.Run("error in all body templates", func(t *testing.T) {
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

	t.Run("error in all target templates", func(t *testing.T) {
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
							Template:    fmt.Sprintf("%s/default.html", tmpTemplateDir),
							ShoutrrrURL: "generic://{{notvalid}}",
						},
					},
				},
			},
		})
		template.CreateDefaultTemplate(tmpTemplateDir)

		app, router := NewApiTest("/f")
		api.PostForm(router)
		r := PerformRequestWithBody(app, "POST", "/f/example", "application/x-www-form-urlencoded", "a=a&b=b")
		assert.Equal(t, http.StatusInternalServerError, r.Code)

		var response api.ResponseBody
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, "Internal server error", response.Message)
		os.RemoveAll(tmpTemplateDir)
	})
}
