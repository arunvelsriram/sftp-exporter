package config

import (
	c "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/spf13/viper"
)

type SFTPConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type Config interface {
	GetBindAddress() string
	GetPort() int
	GetLogLevel() string
	GetSFTPConfig() SFTPConfig
}

type sftpExporterConfig struct {
	BindAddress string
	Port        int
	LogLevel    string
	SFTPConfig  SFTPConfig
}

func LoadConfig() Config {
	return sftpExporterConfig{
		BindAddress: viper.GetString(c.ViperKeyBindAddress),
		Port:        viper.GetInt(c.ViperKeyPort),
		LogLevel:    viper.GetString(c.ViperKeyLogLevel),
		SFTPConfig: SFTPConfig{
			Host: viper.GetString(c.ViperKeySFTPHost),
			Port: viper.GetInt(c.ViperKeySFTPPort),
			User: viper.GetString(c.ViperKeySFTPUser),
			Pass: viper.GetString(c.ViperKeySFTPPass),
		},
	}
}

func (c sftpExporterConfig) GetBindAddress() string {
	return c.BindAddress
}

func (c sftpExporterConfig) GetPort() int {
	return c.Port
}

func (c sftpExporterConfig) GetSFTPConfig() SFTPConfig {
	return c.SFTPConfig
}

func (c sftpExporterConfig) GetLogLevel() string {
	return c.LogLevel
}
