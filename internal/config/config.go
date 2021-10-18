package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// TargetConfig is the config of a shoutrrr target
type TargetConfig struct {
	Enabled     bool
	Template    string
	ShoutrrrURL string
	Params      map[string]interface{}
}

// FormConfig is the config of a form
type FormConfig struct {
	Enabled bool
	Id      string
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
	Forms  []*FormConfig
}

type OnConfigChangeFunc func(event fsnotify.Event, config *Config, oldConfig *Config)

var conf *Config = DefaultConfig()
var onConfigChangeFuncs []OnConfigChangeFunc

// Setup reads the config file and parses it
func Setup(configPath string) error {
	return SetupWithName(configPath, "config")
}

// Setup reads the config file and parses it and allows to set the name
func SetupWithName(configPath string, configName string) error {

	viper.Reset()
	viper.SetEnvPrefix("FORMRECEVR_")
	viper.SetConfigType("yml")
	viper.SetConfigName(configName)
	viper.SetConfigPermissions(0600)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath("./formrecevr-config")

	viper.SetDefault("listen", DefaultConfig().Listen)

	if err := viper.ReadInConfig(); err != nil {

		if _, e := err.(viper.ConfigFileNotFoundError); !e {
			return err
		}

		os.MkdirAll(configPath, os.ModePerm)
		viper.SetConfigFile(fmt.Sprintf("%s/%s", configPath, fmt.Sprintf("%s.yml", configName)))
		viper.SetDefault("forms", DefaultConfig().Forms)
		if err = viper.WriteConfig(); err != nil {
			return err
		}
	}

	viper.OnConfigChange(handleConfigChanged)
	viper.WatchConfig()

	return viper.Unmarshal(&conf)
}

func OnConfigChange(callback OnConfigChangeFunc) {
	onConfigChangeFuncs = append(onConfigChangeFuncs, callback)
}

// GetConfig returns the config
func GetConfig() *Config {
	return conf
}

func ConfigPathUsed() string {
	return path.Dir(viper.ConfigFileUsed())
}

// DefaultConfig returns the default config
func DefaultConfig() *Config {
	return &Config{
		Listen: ListenConfig{
			Host: "0.0.0.0",
			Port: 8088,
		},
		Forms: []*FormConfig{
			{
				Enabled: false,
				Id:      "Example",
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
	if c == nil {
		conf = nil
		viper.Reset()
	} else {
		viper.Set("forms", c.Forms)
		viper.Set("listen", c.Listen)
		conf = c
	}
}

func handleConfigChanged(e fsnotify.Event) {
	log.Println("Config file changed")
	oldConfig := *conf
	conf = nil
	viper.Unmarshal(&conf)
	for _, callbackFunc := range onConfigChangeFuncs {
		callbackFunc(e, conf, &oldConfig)
	}
}
