package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type TargetConfig struct {
	Enabled     bool
	Template    string
	ShoutrrrURL string
}

type FormConfig struct {
	Enabled bool
	Targets []*TargetConfig
}

type ListenConfig struct {
	Host string
	Port int
}

type Config struct {
	Listen ListenConfig
	Forms  map[string]*FormConfig
}

var conf *Config = DefaultConfig()

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
	fmt.Printf("--- t:\n%v\n\n", conf)
	return nil
}

func GetConfig() *Config {
	return conf
}

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
						Template:    "./templates/default.html",
						ShoutrrrURL: "See: https://containrrr.dev/shoutrrr/v0.5/services/overview/",
					},
				},
			},
		},
	}
}
