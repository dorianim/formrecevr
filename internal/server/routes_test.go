package server

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_registerRoutes(t *testing.T) {
	registerRoutes(gin.New())
}
