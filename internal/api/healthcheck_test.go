package api_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dorianim/formrecevr/internal/api"
)

func TestHealthcheckApi(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		app, router := NewApiTest("/api/v1")
		api.Healthcheck(router)
		r := PerformRequest(app, "GET", "/api/v1/healthcheck")
		assert.Equal(t, http.StatusOK, r.Code)
	})
}
