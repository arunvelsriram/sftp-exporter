package config

import (
	"encoding/base64"
	"fmt"

	c "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/arunvelsriram/sftp-exporter/pkg/utils"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type sftpConfig struct {
	Host          string
	Port          int
	User          string
	Pass          string
	Key           []byte
	KeyFile       string
	KeyPassphrase string
}

type Config interface {
	GetBindAddress() string
	GetPort() int
	GetLogLevel() string
	GetSFTPHost() string
	GetSFTPPort() int
	GetSFTPUser() string
	GetSFTPPass() string
	GetSFTPKey() []byte
	GetSFTPKeyFile() string
	GetSFTPKeyPassphrase() string
}

type sftpExporterConfig struct {
	BindAddress string
	Port        int
	LogLevel    string
	SFTPConfig  sftpConfig
}

func resolveKey(encodedKey, keyfile string, fs afero.Fs) (sftpKey []byte, err error) {
	if utils.IsNotEmpty(encodedKey) && utils.IsNotEmpty(keyfile) {
		return sftpKey, fmt.Errorf("only one of key or keyfile should be specified")
	}

	if utils.IsNotEmpty(encodedKey) {
		sftpKey, err = base64.StdEncoding.DecodeString(encodedKey)
		if err != nil {
			return sftpKey, err
		}
	} else if utils.IsNotEmpty(keyfile) {
		sftpKey, err = afero.ReadFile(fs, keyfile)
		if err != nil {
			return sftpKey, err
		}
	}
	return sftpKey, nil
}

func LoadConfig(fs afero.Fs) (Config, error) {
	encodedKey := viper.GetString(c.ViperKeySFTPKey)
	keyFile := viper.GetString(c.ViperKeySFTPKeyFile)
	key, err := resolveKey(encodedKey, keyFile, fs)
	if err != nil {
		return nil, err
	}

	return sftpExporterConfig{
		BindAddress: viper.GetString(c.ViperKeyBindAddress),
		Port:        viper.GetInt(c.ViperKeyPort),
		LogLevel:    viper.GetString(c.ViperKeyLogLevel),
		SFTPConfig: sftpConfig{
			Host:          viper.GetString(c.ViperKeySFTPHost),
			Port:          viper.GetInt(c.ViperKeySFTPPort),
			User:          viper.GetString(c.ViperKeySFTPUser),
			Pass:          viper.GetString(c.ViperKeySFTPPass),
			Key:           key,
			KeyFile:       keyFile,
			KeyPassphrase: viper.GetString(c.ViperKeySFTPKeyPassphrase),
		},
	}, nil
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

func (c sftpExporterConfig) GetSFTPKey() []byte {
	return c.SFTPConfig.Key
}

func (c sftpExporterConfig) GetSFTPKeyFile() string {
	return c.SFTPConfig.KeyFile
}

func (c sftpExporterConfig) GetSFTPKeyPassphrase() string {
	return c.SFTPConfig.KeyPassphrase
}
