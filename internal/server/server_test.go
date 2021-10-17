package server

import (
	"testing"

	"github.com/dorianim/formrecevr/internal/config"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

func TestListenAndServe(t *testing.T) {
	t.Run("invalid port", func(t *testing.T) {
		config.SetConfig(&config.Config{
			Listen: config.ListenConfig{
				Host: "0.0.0.0",
				Port: 999999999,
			},
		})
		Setup()
		err := ListenAndServe()
		assert.NotNil(t, err)
	})
}

func TestHandleConfigChanged(t *testing.T) {
	c1 := config.Config{
		Listen: config.ListenConfig{
			Host: "0.0.0.0",
			Port: 1,
		},
	}
	c2 := config.Config{
		Listen: config.ListenConfig{
			Host: "0.0.0.0",
			Port: 2,
		},
	}
	handleConfigChanged(fsnotify.Event{}, &c1, &c2)
}
