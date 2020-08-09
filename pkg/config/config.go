package config

import (
	c "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/spf13/viper"
)

type sftpConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type Config interface {
	GetBindAddress() string
	GetPort() int
	GetLogLevel() string
	GetSFTPHost() string
	GetSFTPPort() int
	GetSFTPUser() string
	GetSFTPPass() string
}

type sftpExporterConfig struct {
	BindAddress string
	Port        int
	LogLevel    string
	SFTPConfig  sftpConfig
}

func LoadConfig() Config {
	return sftpExporterConfig{
		BindAddress: viper.GetString(c.ViperKeyBindAddress),
		Port:        viper.GetInt(c.ViperKeyPort),
		LogLevel:    viper.GetString(c.ViperKeyLogLevel),
		SFTPConfig: sftpConfig{
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

func (c sftpExporterConfig) GetLogLevel() string {
	return c.LogLevel
}

func (c sftpExporterConfig) GetSFTPHost() string {
	return c.SFTPConfig.Host
}

func (c sftpExporterConfig) GetSFTPPort() int {
	return c.SFTPConfig.Port
}

func (c sftpExporterConfig) GetSFTPUser() string {
	return c.SFTPConfig.User
}

func (c sftpExporterConfig) GetSFTPPass() string {
	return c.SFTPConfig.Pass
}