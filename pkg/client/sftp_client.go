package client

import (
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
	s.sshClient, err = NewSSHClient()
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

func NewSFTPClient() SFTPClient {
	return &sftpClient{}
}
