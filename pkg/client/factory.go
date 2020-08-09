package client

import "github.com/arunvelsriram/sftp-exporter/pkg/config"

type (
	Factory interface {
		SFTPClient() (SFTPClient, error)
	}

	factory struct {
		config config.Config
	}
)

func (p factory) SFTPClient() (SFTPClient, error) {
	return newSFTPClient(p.config)
}

func NewFactory(c config.Config) Factory {
	return factory{c}
}
