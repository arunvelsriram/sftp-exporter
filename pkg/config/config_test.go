package config

import (
	"testing"

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

func (s *ConfigTestSuite) TestConfigLoadConfig() {
	s.Run("should return config", func() {
		actual, err := LoadConfig(s.fs)

		s.NoError(err)
		s.Equal(sftpExporterConfig{}, actual)
	})

	s.Run("should return err when both key and keyfile are provided", func() {
		viper.Reset()
		viper.Set(c.ViperKeySFTPKey, "sftp private key")
		viper.Set(c.ViperKeySFTPKeyFile, "sftp private keyfile")

		_, err := LoadConfig(s.fs)

		s.EqualError(err, "only one of key or keyfile should be specified")
	})

	s.Run("should return error if key is not encoded properly", func() {
		viper.Reset()
		viper.Set(c.ViperKeySFTPKey, "invalid encoding")

		_, err := LoadConfig(s.fs)

		s.EqualError(err, "illegal base64 data at input byte 7")
	})

	s.Run("should store decoded key for given encoded key", func() {
		viper.Reset()
		viper.Set(c.ViperKeySFTPKey, "YXJ1bg==")

		actual, err := LoadConfig(s.fs)

		s.NoError(err)
		s.Equal([]byte("arun"), actual.GetSFTPKey())
	})

	s.Run("should return error if reading keyfile fails", func() {
		viper.Reset()
		viper.Set(c.ViperKeySFTPKeyFile, "invalidfile")

		_, err := LoadConfig(s.fs)

		s.EqualError(err, "open invalidfile: file does not exist")
	})

	s.Run("should store key and keyfile for a valid keyfile", func() {
		file, _ := afero.TempFile(s.fs, "", "config-test")
		_, _ = file.WriteString("private key")
		viper.Reset()
		viper.Set(c.ViperKeySFTPKeyFile, file.Name())

		actual, err := LoadConfig(s.fs)

		s.NoError(err)
		s.Equal(file.Name(), actual.GetSFTPKeyFile())
		s.Equal([]byte("private key"), actual.GetSFTPKey())
	})
}
