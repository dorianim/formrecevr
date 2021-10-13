package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/accounts
func Healthcheck(router *gin.RouterGroup) {
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})
}
