package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigmapGet(t *testing.T) {
	t.Run("when key is present", func(t *testing.T) {
		c := Configmap{"key": "value"}

		value, ok := c.Get("key")

		assert.True(t, ok)
		assert.Equal(t, "value", value)
	})

	t.Run("when key is not present", func(t *testing.T) {
		c := Configmap{"key": "value"}

		_, ok := c.Get("invalid")

		assert.False(t, ok)
	})
}

func TestConfigmapSet(t *testing.T) {
	c := Configmap{"key": "value"}

	c.Set("key", "newvalue")

	expected := Configmap{"key": "newvalue"}
	assert.Equal(t, expected, c)
}
