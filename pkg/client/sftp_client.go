package client

import (
	"fmt"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/arunvelsriram/sftp-exporter/pkg/utils"
	"github.com/kr/fs"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type (
	SFTPClient interface {
		Connect() error
		Close() error
		StatVFS(path string) (*sftp.StatVFS, error)
		Walk(root string) *fs.Walker
	}

	sftpClient struct {
		*sftp.Client
		sshClient *ssh.Client
		config    config.Config
	}
)

func (s *sftpClient) Close() error {
	if err := s.Client.Close(); err != nil {
		log.WithFields(log.Fields{
			"event": "closing SFTP connection"},
		).Error(err)
	}
	if err := s.sshClient.Close(); err != nil {
		log.WithFields(log.Fields{
			"event": "closing SSH connection"},
		).Error(err)
	}
	return nil
}

func (s *sftpClient) Connect() (err error) {
	addr := fmt.Sprintf("%s:%d", s.config.GetSFTPHost(), s.config.GetSFTPPort())
	auth, err := utils.SSHAuthMethods(s.config.GetSFTPPass(), s.config.GetSFTPKey(), s.config.GetSFTPKeyPassphrase())
	if err != nil {
		log.Error("unable to get SSH auth methods")
		return err
	}
	clientConfig := &ssh.ClientConfig{
		User:            s.config.GetSFTPUser(),
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	s.sshClient, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return err
	}

	s.Client, err = sftp.NewClient(s.sshClient)
	if err != nil {
		if err := s.sshClient.Close(); err != nil {
			log.WithFields(log.Fields{
				"event": "closing SSH connection"},
			).Error(err)
		}
		return err
	}
	return nil
}

func NewSFTPClient(cfg config.Config) SFTPClient {
	return &sftpClient{config: cfg}
}
