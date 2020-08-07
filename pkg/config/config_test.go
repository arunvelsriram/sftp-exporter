package config

import (
	"testing"

	. "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	viper.Set(ViperKeyPort, 8080)
	viper.Set(ViperKeySFTPHost, "localhost")
	viper.Set(ViperKeySFTPPort, 22)
	viper.Set(ViperKeySFTPUser, "arun")
	viper.Set(ViperKeySFTPPass, "arun@123")

	c := NewConfig()

	expected := sftpExporterConfig{
		Port: 8080,
		SFTPConfig: SFTPConfig{
			Host: "localhost",
			Port: 22,
			User: "arun",
			Pass: "arun@123",
		},
	}
	assert.Equal(t, expected, c)
}
