package config_test

import (
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

	c := config.GetConfig()
	assert.IsType(t, new(config.Config), c)
	assert.IsType(t, new(map[string]*config.FormConfig), &c.Forms)

	assert.Equal(t, c.Listen.Host, "0.0.0.0")
	assert.Equal(t, c.Listen.Port, 8088)

	form1id := "6db3fc2c-9220-4fd1-8e40-c900c060ca9e"
	assert.IsType(t, new(config.FormConfig), c.Forms[form1id])

	assert.True(t, c.Forms[form1id].Enabled)
	assert.IsType(t, new([]*config.TargetConfig), &c.Forms[form1id].Targets)

	targets1 := c.Forms[form1id].Targets
	assert.Len(t, targets1, 2)
	assert.False(t, targets1[0].Enabled)
	assert.Equal(t, targets1[0].Template, "./templates/default.html")
	assert.Equal(t, targets1[0].ShoutrrrURL, "invalid://invalid")

	assert.True(t, targets1[1].Enabled)
	assert.Equal(t, targets1[1].Template, "./templates/default.html")
	assert.Equal(t, targets1[1].ShoutrrrURL, "telegram://<token>@telegram?channels=<channel>&ParseMode=none")
}
