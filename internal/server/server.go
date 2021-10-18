package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/dorianim/formrecevr/internal/config"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
)

var httpServer *http.Server
var httpRouter *gin.Engine

// Setup sets up dan configures the server
func Setup() {
	httpRouter = gin.New()

	httpRouter.Use(gin.Logger())
	registerRoutes(httpRouter)

	config.OnConfigChange(handleConfigChanged)
}

// ListenAndServe starts the server. It is a blocking function.
func ListenAndServe() error {
	var err error = nil
	for err == nil || err == http.ErrServerClosed {
		c := config.GetConfig()
		httpServer = &http.Server{
			Handler: httpRouter,
			Addr:    fmt.Sprintf("%s:%d", c.Listen.Host, c.Listen.Port),
		}

		log.Printf("Starting web server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
	}
	return err
}

func handleConfigChanged(event fsnotify.Event, config *config.Config, oldConfig *config.Config) {
	if !reflect.DeepEqual(oldConfig.Listen, config.Listen) {
		log.Println("Listen config changed, restarting web server ...")
		httpServer.Shutdown(context.Background())
	}
}
