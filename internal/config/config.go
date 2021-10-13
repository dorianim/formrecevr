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

func Init() {
	if conf != nil {
		return
	}

	content, err := ioutil.ReadFile("./config.yml")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	} else if os.IsNotExist(err) {
		log.Println("Config empty")
	}

	// Convert []byte to string and print to screen
	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", conf)
}

func GetConfig() *Config {
	return conf
}
