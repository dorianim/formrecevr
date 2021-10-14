package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

type Config struct {
	Listen struct {
		Host string
		Port int
	}

	Forms map[string]*FormConfig
}

var conf *Config = nil

func NewConfig(configFilePath string) error {
	if conf != nil {
		log.Printf("Config already set")
		return nil
	}

	content, err := ioutil.ReadFile(configFilePath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
		return err
	} else if os.IsNotExist(err) {
		log.Printf("Config %s empty", configFilePath)
	}

	// Convert []byte to string and print to screen
	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
		return err
	}
	fmt.Printf("--- t:\n%v\n\n", conf)
	return nil
}

func GetConfig() *Config {
	return conf
}
