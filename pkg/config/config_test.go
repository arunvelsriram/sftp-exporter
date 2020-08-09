package config

import (
	"testing"

	. "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	viper.Set(ViperKeyBindAddress, "127.0.0.1")
	viper.Set(ViperKeyPort, 8080)
	viper.Set(ViperKeySFTPHost, "localhost")
	viper.Set(ViperKeySFTPPort, 22)
	viper.Set(ViperKeySFTPUser, "arun")
	viper.Set(ViperKeySFTPPass, "arun@123")

	c := LoadConfig()

	expected := sftpExporterConfig{
		BindAddress: "127.0.0.1",
		Port:        8080,
		SFTPConfig: sftpConfig{
			Host: "localhost",
			Port: 22,
			User: "arun",
			Pass: "arun@123",
		},
	}
	assert.Equal(t, expected, c)
}
