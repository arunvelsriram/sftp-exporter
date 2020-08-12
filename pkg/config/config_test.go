package config_test

import (
	"testing"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	c "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	fs afero.Fs
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) SetupTest() {
	s.fs = afero.NewMemMapFs()
}

func (s *ConfigTestSuite) TearDownTest() {
	viper.Reset()
}

func (s *ConfigTestSuite) TestConfigLoadConfig() {
	viper.Set(c.ViperKeySFTPUser, "sftp user")
	viper.Set(c.ViperKeySFTPPass, "sftp pass")
	actual := config.MustLoadConfig(s.fs)

	s.Equal("sftp user", actual.GetSFTPUser())
	s.Equal("sftp pass", actual.GetSFTPPass())
}

func (s *ConfigTestSuite) TestConfigLoadConfigPanicsIfSFTPEmptyUser() {
	s.PanicsWithValue("config sftp_user is required", func() { config.MustLoadConfig(s.fs) })
}

func (s *ConfigTestSuite) TestConfigLoadConfigPanicsWhenBothAllAuthTypesAreEmpty() {
	viper.Set(c.ViperKeySFTPUser, "sftp user")

	s.PanicsWithValue("either one of sftp_pass, sftp_key, sftp_key_file is required", func() { config.MustLoadConfig(s.fs) })
}

func (s *ConfigTestSuite) TestConfigLoadConfigPanicsWhenBothKeyAndKeyFileAreGiven() {
	viper.Set(c.ViperKeySFTPUser, "sftp user")
	viper.Set(c.ViperKeySFTPKey, "sftp private key")
	viper.Set(c.ViperKeySFTPKeyFile, "sftp private keyfile")

	s.PanicsWithValue("only one of sftp_key, sftp_key_file should be provided", func() { config.MustLoadConfig(s.fs) })
}
func (s *ConfigTestSuite) TestConfigLoadConfigPanicsForInvalidKey() {
	viper.Set(c.ViperKeySFTPUser, "sftp user")
	viper.Set(c.ViperKeySFTPKey, "invalid encoding")

	s.PanicsWithValue("illegal base64 data at input byte 7", func() { config.MustLoadConfig(s.fs) })
}

func (s *ConfigTestSuite) TestConfigLoadConfigStoresDecodedKey() {
	viper.Set(c.ViperKeySFTPUser, "sftp user")
	viper.Set(c.ViperKeySFTPKey, "YXJ1bg==")

	actual := config.MustLoadConfig(s.fs)

	s.Equal([]byte("arun"), actual.GetSFTPKey())
}

func (s *ConfigTestSuite) TestConfigLoadConfigPanicsIfReadingKeyFileFails() {
	viper.Set(c.ViperKeySFTPUser, "sftp user")
	viper.Set(c.ViperKeySFTPKeyFile, "invalidfile")

	s.PanicsWithValue("open invalidfile: file does not exist", func() { config.MustLoadConfig(s.fs) })
}

func (s *ConfigTestSuite) TestConfigLoadConfigStoresKeyForValidKeyFile() {
	file, _ := afero.TempFile(s.fs, "", "config-test")
	_, _ = file.WriteString("private key")
	viper.Set(c.ViperKeySFTPUser, "sftp user")
	viper.Set(c.ViperKeySFTPKeyFile, file.Name())

	actual := config.MustLoadConfig(s.fs)

	s.Equal(file.Name(), actual.GetSFTPKeyFile())
	s.Equal([]byte("private key"), actual.GetSFTPKey())
}
