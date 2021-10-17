package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// TargetConfig is the config of a shoutrrr target
type TargetConfig struct {
	Enabled     bool
	Template    string
	ShoutrrrURL string
}

// FormConfig is the config of a form
type FormConfig struct {
	Enabled bool
	Targets []*TargetConfig
}

// ListenConfig is for the server
type ListenConfig struct {
	Host string
	Port int
}

// Config is the type which contains all above sub-configs
type Config struct {
	Listen ListenConfig
	Forms  map[string]*FormConfig
}

var conf *Config = DefaultConfig()

// NewConfig reads a configfile, parses it and puts it into conf
func NewConfig(configFilePath string) error {
	conf = nil
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if os.IsNotExist(err) {
		log.Printf("Config %s empty, writing default to file", configFilePath)
		conf = DefaultConfig()
		return WriteConfigToFile(configFilePath, conf)
	}

	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		return err
	}

	return nil
}

// GetConfig returns the config
func GetConfig() *Config {
	return conf
}

// WriteConfigToFile writes a config to a file
func WriteConfigToFile(configFilePath string, conf *Config) error {
	d, _ := yaml.Marshal(&conf)
	os.MkdirAll(filepath.Dir(configFilePath), os.ModePerm)

	f, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	f.Truncate(0)
	f.Write(d)

	f.Close()
	return nil
}

// DefaultConfig returns the default config
func DefaultConfig() *Config {
	return &Config{
		Listen: ListenConfig{
			Host: "0.0.0.0",
			Port: 8088,
		},
		Forms: map[string]*FormConfig{
			"Example": {
				Enabled: false,
				Targets: []*TargetConfig{
					{
						Enabled:     false,
						Template:    "default.html",
						ShoutrrrURL: "See: https://containrrr.dev/shoutrrr/v0.5/services/overview/",
					},
				},
			},
		},
	}
}

// SetConfig sets the config to c
func SetConfig(c *Config) {
	conf = c
}
