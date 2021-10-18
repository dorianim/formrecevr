package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		e := Setup("../../testdata")
		assert.Nil(t, e)
	})

	t.Run("invalid config", func(t *testing.T) {
		e := SetupWithName("../../testdata", "invalid-config")
		assert.NotNil(t, e)
	})

	t.Run("non existent config", func(t *testing.T) {
		tmpPath := "../../testdata/non-existent"
		e := os.RemoveAll(tmpPath)
		assert.True(t, e == nil || os.IsNotExist(e))

		// check if default config works
		e = Setup(tmpPath)
		assert.Nil(t, e)
		assert.Equal(t, DefaultConfig(), GetConfig())

		// check if config was created
		_, err := os.Stat(fmt.Sprintf("%s/config.yml", tmpPath))
		assert.False(t, os.IsNotExist(err))

		assert.Nil(t, os.RemoveAll(tmpPath))
	})

	t.Run("invalid path", func(t *testing.T) {
		nullByte := []byte{0x0}
		tmpPath := string(nullByte)

		e := Setup(tmpPath)
		assert.NotNil(t, e)
	})
}

func TestGetConfig(t *testing.T) {
	e := Setup("../../testdata")
	log.Println(e)
	assert.Nil(t, e)

	assert.Equal(t, DefaultConfig(), GetConfig())
}

func TestSetConfig(t *testing.T) {
	SetConfig(nil)
	SetConfig(DefaultConfig())
	assert.Equal(t, DefaultConfig(), GetConfig())
}

func TestHandleConfigChange(t *testing.T) {
	t.Run("parameter changed", func(t *testing.T) {
		ok := false

		OnConfigChange(func(event fsnotify.Event, config *Config, oldConfig *Config) {
			ok = true
		})

		handleConfigChanged(fsnotify.Event{})
		assert.True(t, ok)
	})
}

func TestConfigPathUsed(t *testing.T) {
	Setup("../../testdata")
	path, _ := filepath.Abs("../../testdata")
	assert.Equal(t, path, PathUsed())
}
