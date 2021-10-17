package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dorianim/formrecevr/internal/config"
	"github.com/dorianim/formrecevr/internal/server"
	"github.com/dorianim/formrecevr/internal/template"
	"github.com/gin-gonic/gin"

	"github.com/spf13/cobra"
)

var rootCmd = NewRootCommand()

// NewRootCommand creates the root command for watchtower
func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "formrecevr",
		Short: "Receives form data via an http API and forwards them",
		Long: `
	Formrecevr receives form data via an http API and forwards them using shoutrrr
	More information is available at: https://github.com/dorianim/formrecevr
	`,
		Run:    Run,
		PreRun: PreRun,
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// PreRun is currently unused
func PreRun(cmd *cobra.Command, _ []string) {
}

// Run starts the server
func Run(c *cobra.Command, names []string) {
	gin.SetMode(gin.ReleaseMode)

	configPath, isSet := os.LookupEnv("FORMRECEVR_CONFIG_PATH")
	if !isSet {
		configPath = "/config"
	}

	configFilePath := fmt.Sprintf("%s/config.yml", configPath)
	if err := config.NewConfig(configFilePath); err != nil {
		log.Fatalf("Error reading config.yml: %v", err)
	}

	templatePath := fmt.Sprintf("%s/templates", configPath)
	template.Setup(templatePath)
	if err := template.CreateDefaultTemplate(); err != nil {
		log.Fatalf("Error creating default template: %v", err)
	}

	s := server.New()
	if err := s.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("http: web server shutdown complete")
		} else {
			log.Fatalf("http: web server closed unexpect: %s", err)
		}
	}
}
