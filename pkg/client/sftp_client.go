package client

import (
	"fmt"
	log "github.com/sirupsen/logrus"

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
	addr := fmt.Sprintf("%s:%d", cfg.GetSFTPHost(), cfg.GetSFTPPort())
	clientConfig := &ssh.ClientConfig{
		User: cfg.GetSFTPUser(),
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.GetSFTPPass()),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		if err := sshClient.Close(); err != nil {
			log.WithFields(log.Fields{
				"event": "closing SFTP connection"},
			).Error(err)
		}
		return nil, err
	}

	return defaultSFTPClient{
		sshClient:  sshClient,
		sftpClient: sftpClient,
		config:     cfg,
	}, nil
}

func (d defaultSFTPClient) Close() {
	if err := d.sftpClient.Close(); err != nil {
		log.WithFields(log.Fields{
			"event": "closing SFTP connection"},
		).Error(err)
	}
	if err := d.sshClient.Close(); err != nil {
		log.WithFields(log.Fields{
			"event": "closing SSH connection"},
		).Error(err)
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
