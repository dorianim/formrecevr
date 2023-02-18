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

// HCaptchaConfig is the config for HCapthca
type HCaptchaConfig struct {
	Enabled    bool
	PrivateKey string
	Score      float32
}

// FormConfig is the config of a form
type FormConfig struct {
	Enabled  bool
	ID       string
	HCaptcha HCaptchaConfig
	Targets  []*TargetConfig
}

// ListenConfig is for the server
type ListenConfig struct {
	Host                string
	Port                int
	UseForwardedHeaders bool
}

// Config is the type which contains all above sub-configs
type Config struct {
	Listen ListenConfig
	Forms  []*FormConfig
}

// OnConfigChangeFunc is the type of a function which can be used as a callback
type OnConfigChangeFunc func(event fsnotify.Event, config *Config, oldConfig *Config)

var conf *Config = DefaultConfig()
var onConfigChangeFuncs []OnConfigChangeFunc

// Setup reads the config file and parses it
func Setup(configPath string) error {
	return SetupWithName(configPath, "config")
}

// SetupWithName reads the config file and parses it and allows to set the name
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

// OnConfigChange can be used to attach a callback function which is called whenever the config changes
func OnConfigChange(callback OnConfigChangeFunc) {
	onConfigChangeFuncs = append(onConfigChangeFuncs, callback)
}

// GetConfig returns the config
func GetConfig() *Config {
	return conf
}

// PathUsed returns the used config path
func PathUsed() string {
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
				ID:      "Example",
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
