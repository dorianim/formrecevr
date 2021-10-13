package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dorianim/formrecevr/internal/config"
	"github.com/gin-gonic/gin"
)

func Start() {
	config := config.GetConfig()
	router := gin.New()

	registerRoutes(router)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Listen.Host, config.Listen.Port),
		Handler: router,
	}

	log.Printf("http: starting web server at %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("http: web server shutdown complete")
		} else {
			log.Fatalf("http: web server closed unexpect: %s", err)
		}
	}
}
