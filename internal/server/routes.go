package server

import (
	"github.com/gin-gonic/gin"

	"github.com/dorianim/formrecevr/internal/api"
)

func registerRoutes(router *gin.Engine) {
	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	router.RedirectTrailingSlash = true

	forms := router.Group("/f")
	{
		api.PostForm(forms)
	}

	apiV1 := router.Group("/api/v1")
	{
		api.Healthcheck(apiV1)
	}
}
