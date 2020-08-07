package config

import (
	c "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/spf13/viper"
)

type Config interface {
	GetPort() int
	GetSFTPConfig() SFTPConfig
}

type sftpExporterConfig struct {
	Port       int
	SFTPConfig SFTPConfig
}

func LoadConfig() Config {
	return sftpExporterConfig{
		Port: viper.GetInt(c.ViperKeyPort),
		SFTPConfig: SFTPConfig{
			Host: viper.GetString("sftp_host"),
			Port: viper.GetInt("sftp_port"),
			User: viper.GetString("sftp_user"),
			Pass: viper.GetString("sftp_pass"),
		},
	}
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
