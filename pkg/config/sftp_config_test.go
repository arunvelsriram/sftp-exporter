package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSFTPConfigToConfigmap(t *testing.T) {
	c := SFTPConfig{
		Host: "localhost",
		Port: 22,
		User: "arun",
		Pass: "arun@123",
	}

	actual := c.toConfigMap()

	expected := Configmap{
		"host": "localhost",
		"port": "22",
		"user": "arun",
		"pass": "arun@123",
	}
	assert.Equal(t, expected, actual)
}
