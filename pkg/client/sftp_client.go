package client

import (
	"fmt"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/arunvelsriram/sftp-exporter/pkg/model"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTPClient interface {
	Close()
	FSStat() (*model.FSStat, error)
}

type defaultSFTPClient struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
	config     config.Config
}

func NewSFTPClient(cfg config.Config) (SFTPClient, error) {
	sftpConfig := cfg.GetSFTPConfig()
	addr := fmt.Sprintf("%s:%d", sftpConfig.Host, sftpConfig.Port)
	clientConfig := &ssh.ClientConfig{
		User: sftpConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(sftpConfig.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	return defaultSFTPClient{
		sshClient:  sshClient,
		sftpClient: sftpClient,
		config:     cfg,
	}, nil
}

func (d defaultSFTPClient) Close() {
	sftpErr := d.sftpClient.Close()
	sshErr := d.sshClient.Close()
	if sftpErr != nil || sshErr != nil {
		fmt.Printf("failed to close connections\nSFTP: %v\nSSH: %v", sftpErr, sshErr)
	}
}

func (d defaultSFTPClient) FSStat() (*model.FSStat, error) {
	statVFS, err := d.sftpClient.StatVFS("/upload")
	if err != nil {
		return nil, err
	}

	return &model.FSStat{
		TotalSpace: float64(statVFS.TotalSpace()),
		FreeSpace:  float64(statVFS.FreeSpace()),
	}, nil
}
