package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/dorianim/formrecevr/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		e := config.NewConfig("../../testdata/config.yml")
		assert.Nil(t, e)
	})

	t.Run("invalid config", func(t *testing.T) {
		e := config.NewConfig("../../testdata/invalid-config.yml")
		assert.NotNil(t, e)
	})

	t.Run("unreadable config", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile("", "unreadable-config.yml")
		assert.Nil(t, err)
		err = os.Chmod(tmpfile.Name(), 0000)
		assert.Nil(t, err)

		e := config.NewConfig(tmpfile.Name())
		assert.NotNil(t, e)
	})

	t.Run("non existent config", func(t *testing.T) {
		tmpFile := "../../testdata/non-existent-config.yml"
		e := os.Remove(tmpFile)
		assert.True(t, e == nil || os.IsNotExist(e))

		// check if default config works
		e = config.NewConfig(tmpFile)
		assert.Nil(t, e)
		assert.Equal(t, config.DefaultConfig(), config.GetConfig())

		// check if default config has been written correctly
		e = config.NewConfig(tmpFile)
		assert.Nil(t, e)
		assert.Equal(t, config.DefaultConfig(), config.GetConfig())

		assert.Nil(t, os.Remove(tmpFile))
	})
}

func TestGetConfig(t *testing.T) {
	e := config.NewConfig("../../testdata/config.yml")
	assert.Nil(t, e)

	assert.Equal(t, config.DefaultConfig(), config.GetConfig())
}

func TestWriteConfigToFile(t *testing.T) {
	t.Run("wrong file permissions", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile("", "unreadable-config.yml")
		assert.Nil(t, err)
		err = os.Chmod(tmpfile.Name(), 0000)
		assert.Nil(t, err)

		err = config.WriteConfigToFile(tmpfile.Name(), nil)
		assert.NotNil(t, err)
	})
}

func TestSetConfig(t *testing.T) {
	config.SetConfig(nil)
	config.SetConfig(config.DefaultConfig())
	assert.Equal(t, config.DefaultConfig(), config.GetConfig())
}
