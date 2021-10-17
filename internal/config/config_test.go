package config

import (
	"io/ioutil"
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

	t.Run("unreadable config", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile("", "unreadable-config.yml")
		assert.Nil(t, err)
		err = os.Chmod(tmpfile.Name(), 0000)
		assert.Nil(t, err)

		e := SetupWithName(filepath.Base(tmpfile.Name()), "unreadable-config")
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

		// check if default config has been written correctly
		e = Setup(tmpPath)
		assert.Nil(t, e)
		assert.Equal(t, DefaultConfig(), GetConfig())

		assert.Nil(t, os.RemoveAll(tmpPath))
	})
}

func TestGetConfig(t *testing.T) {
	e := Setup("../../testdata/config.yml")
	assert.Nil(t, e)

	assert.Equal(t, DefaultConfig(), GetConfig())
}

func TestSetConfig(t *testing.T) {
	SetConfig(nil)
	SetConfig(DefaultConfig())
	assert.Equal(t, DefaultConfig(), GetConfig())
}

func TestHandleConfigChange(t *testing.T) {
	ok := false
	OnConfigChange(func(event fsnotify.Event, config *Config, oldConfig *Config) {
		ok = true
	})
	handleConfigChanged(fsnotify.Event{})
	assert.True(t, ok)
}
