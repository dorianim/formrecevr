package main

import (
	"log"
	"os"

	"github.com/dorianim/formrecevr/internal/config"
	"github.com/dorianim/formrecevr/internal/server"

	"github.com/spf13/cobra"
)

var version = "development"

//var log = event.Log

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
	configFilePath := os.Getenv("FORMRECEVR_CONFIG_FILE_PATH")
	if configFilePath == "" {
		configFilePath = "/config/config.yml"
	}

	if err := config.NewConfig(configFilePath); err != nil {
		log.Fatalf("Error reading config.yml: %v", err)
	}

	server.Start()
}
