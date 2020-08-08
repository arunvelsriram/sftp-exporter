package config

import (
	c "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/spf13/viper"
)

type Config interface {
	GetBindAddress() string
	GetPort() int
	GetSFTPConfig() SFTPConfig
}

type sftpExporterConfig struct {
	BindAddress string
	Port        int
	SFTPConfig  SFTPConfig
}

func LoadConfig() Config {
	return sftpExporterConfig{
		BindAddress: viper.GetString(c.ViperKeyBindAddress),
		Port:        viper.GetInt(c.ViperKeyPort),
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

func (c sftpExporterConfig) GetSFTPConfigmap() Configmap {
	return c.SFTPConfig.toConfigMap()
}
