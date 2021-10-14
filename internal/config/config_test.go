package config_test

import (
	"testing"

	"github.com/dorianim/formrecevr/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	e := config.NewConfig("./config.yml")
	assert.Nil(t, e)
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
