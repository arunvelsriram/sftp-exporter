package client

import (
	"fmt"
	"github.com/arunvelsriram/sftp-exporter/pkg/utils"

	log "github.com/sirupsen/logrus"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/arunvelsriram/sftp-exporter/pkg/model"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type (
	SFTPClient interface {
		Close()
		FSStat() (*model.FSStat, error)
	}

	defaultSFTPClient struct {
		sshClient  *ssh.Client
		sftpClient *sftp.Client
		config     config.Config
	}
)

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

func parsePrivateKey(key, keyPassphrase []byte) (parsedKey ssh.Signer, err error) {
	if len(keyPassphrase) > 0 {
		log.Debug("key has passphrase")
		parsedKey, err = ssh.ParsePrivateKeyWithPassphrase(key, keyPassphrase)
		if err != nil {
			log.Error("failed to parse key with passphrase")
			return nil, err
		}
		return parsedKey, err
	}

	log.Debug("key has no passphrase")
	parsedKey, err = ssh.ParsePrivateKey(key)
	if err != nil {
		log.Error("failed to parse key")
		return nil, err
	}
	return parsedKey, err
}

func sshAuthMethods(cfg config.Config) ([]ssh.AuthMethod, error) {
	pass := cfg.GetSFTPPass()
	key := cfg.GetSFTPKey()
	keyPassphrase := cfg.GetSFTPKeyPassphrase()

	if len(key) > 0 && utils.IsNotEmpty(pass) {
		log.Debug("will be authenticating using key and password")
		parsedKey, err := parsePrivateKey(key, keyPassphrase)
		if err != nil {
			return nil, err
		}
		return []ssh.AuthMethod{
			ssh.PublicKeys(parsedKey),
			ssh.Password(pass),
		}, nil
	} else if len(key) > 0 {
		log.Debug("will be authenticating using key")
		parsedKey, err := parsePrivateKey(key, keyPassphrase)
		if err != nil {
			return nil, err
		}
		return []ssh.AuthMethod{
			ssh.PublicKeys(parsedKey),
		}, nil
	} else if utils.IsNotEmpty(pass) {
		log.Debug("will be authenticating using password")
		return []ssh.AuthMethod{
			ssh.Password(pass),
		}, nil
	}

	return nil, fmt.Errorf("either one of password or key is required")
}

func newSFTPClient(cfg config.Config) (SFTPClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.GetSFTPHost(), cfg.GetSFTPPort())
	auth, err := sshAuthMethods(cfg)
	if err != nil {
		log.Error("unable to get SSH auth methods")
		return nil, err
	}
	clientConfig := &ssh.ClientConfig{
		User:            cfg.GetSFTPUser(),
		Auth:            auth,
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
