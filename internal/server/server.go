package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dorianim/formrecevr/internal/config"
	"github.com/gin-gonic/gin"
)

// New starts the server
func New() *http.Server {
	config := config.GetConfig()
	router := gin.New()

	router.Use(gin.Logger())
	registerRoutes(router)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Listen.Host, config.Listen.Port),
		Handler: router,
	}

	log.Printf("http: starting web server at %s", server.Addr)

	return server
}
