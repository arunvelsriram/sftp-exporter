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
	Paths         []string
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
	GetSFTPKeyPassphrase() []byte
	GetSFTPPaths() []string
}

type sftpExporterConfig struct {
	BindAddress string
	Port        int
	LogLevel    string
	SFTPConfig  sftpConfig
}

func mustResolveKey(encodedKey, keyFile string, fs afero.Fs) []byte {
	mustValidateKeySource(encodedKey, keyFile)

	var sftpKey []byte
	var err error
	if utils.IsNotEmpty(encodedKey) {
		sftpKey, err = base64.StdEncoding.DecodeString(encodedKey)
		utils.PanicIfErr(err)
	} else if utils.IsNotEmpty(keyFile) {
		sftpKey, err = afero.ReadFile(fs, keyFile)
		utils.PanicIfErr(err)
	}
	return sftpKey
}

func mustGetString(k string) string {
	v := viper.GetString(k)
	if utils.IsEmpty(v) {
		errMsg := fmt.Sprintf("config %s is required", k)
		panic(errMsg)
	}
	return v
}

func mustValidateKeySource(key, keyFile string) {
	if utils.IsNotEmpty(key) && utils.IsNotEmpty(keyFile) {
		errMsg := fmt.Sprintf("only one of %s, %s should be provided", c.ViperKeySFTPKey, c.ViperKeySFTPKeyFile)
		panic(errMsg)
	}
}

func mustValidateAuthTypes(pass, key, keyFile string) {
	if utils.IsEmpty(pass) && utils.IsEmpty(key) && utils.IsEmpty(keyFile) {
		errMsg := fmt.Sprintf("either one of %s, %s, %s is required", c.ViperKeySFTPPass, c.ViperKeySFTPKey, c.ViperKeySFTPKeyFile)
		panic(errMsg)
	}
}

func MustLoadConfig(fs afero.Fs) Config {
	user := mustGetString(c.ViperKeySFTPUser)

	pass := viper.GetString(c.ViperKeySFTPPass)
	encodedKey := viper.GetString(c.ViperKeySFTPKey)
	keyFile := viper.GetString(c.ViperKeySFTPKeyFile)
	mustValidateAuthTypes(pass, encodedKey, keyFile)

	key := mustResolveKey(encodedKey, keyFile, fs)

	return sftpExporterConfig{
		BindAddress: viper.GetString(c.ViperKeyBindAddress),
		Port:        viper.GetInt(c.ViperKeyPort),
		LogLevel:    viper.GetString(c.ViperKeyLogLevel),
		SFTPConfig: sftpConfig{
			Host:          viper.GetString(c.ViperKeySFTPHost),
			Port:          viper.GetInt(c.ViperKeySFTPPort),
			User:          user,
			Pass:          pass,
			Key:           key,
			KeyFile:       keyFile,
			KeyPassphrase: viper.GetString(c.ViperKeySFTPKeyPassphrase),
			Paths:         viper.GetStringSlice(c.ViperKeySFTPPaths),
		},
	}
}

func NewConfig() Config {
	return sftpExporterConfig{}
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

func (c sftpExporterConfig) GetSFTPKeyPassphrase() []byte {
	return []byte(c.SFTPConfig.KeyPassphrase)
}

func (c sftpExporterConfig) GetSFTPPaths() []string {
	return c.SFTPConfig.Paths
}
